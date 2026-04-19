package handler

import (
	"context"
	"net/http"
	"sync"
	"time"

	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	orderv1 "github.com/Andrew1996-la/ship-builder/shared/pkg/openapi/order/v1"
	inventoryv1 "github.com/Andrew1996-la/ship-builder/shared/pkg/proto/inventory/v1"
	paymentv1 "github.com/Andrew1996-la/ship-builder/shared/pkg/proto/payment/v1"
)

func mapPaymentMethod(pm orderv1.PaymentMethod) (paymentv1.PaymentMethod, error) {
	switch pm {
	case orderv1.PaymentMethodCARD:
		return paymentv1.PaymentMethod_PAYMENT_METHOD_CARD, nil
	case orderv1.PaymentMethodSBP:
		return paymentv1.PaymentMethod_PAYMENT_METHOD_SBP, nil
	case orderv1.PaymentMethodCREDITCARD:
		return paymentv1.PaymentMethod_PAYMENT_METHOD_CREDIT_CARD, nil
	case orderv1.PaymentMethodINVESTORMONEY:
		return paymentv1.PaymentMethod_PAYMENT_METHOD_INVESTOR_MONEY, nil
	default:
		return 0, status.Error(codes.InvalidArgument, "не поддерживаемый payment_method")
	}
}

// Order представляет заказ на постройку космического корабля.
type Order struct {
	OrderUUID       uuid.UUID
	HullUUID        uuid.UUID
	EngineUUID      uuid.UUID
	ShieldUUID      *uuid.UUID // опциональный
	WeaponUUID      *uuid.UUID // опциональный
	TotalPrice      int64      // в копейках
	TransactionUUID *uuid.UUID
	PaymentMethod   *string
	Status          string // PENDING_PAYMENT, PAID, CANCELLED
	CreatedAt       time.Time
}

// OrderStore — хранилище заказов (in-memory).
type OrderStore struct {
	mu     sync.RWMutex
	orders map[uuid.UUID]Order
}

// NewOrderStore создаёт новое пустое хранилище заказов.
func NewOrderStore() *OrderStore {
	return &OrderStore{
		orders: make(map[uuid.UUID]Order),
	}
}

// OrderHandler реализует интерфейс orderv1.Handler, сгенерированный ogen.
type OrderHandler struct {
	orderv1.UnimplementedHandler
	inventoryClient inventoryv1.InventoryServiceClient
	paymentClient   paymentv1.PaymentServiceClient
	store           *OrderStore
}

// NewOrderHandler создаёт новый обработчик заказов.
func NewOrderHandler(
	inventoryClient inventoryv1.InventoryServiceClient,
	paymentClient paymentv1.PaymentServiceClient,
	store *OrderStore,
) *OrderHandler {
	return &OrderHandler{
		inventoryClient: inventoryClient,
		paymentClient:   paymentClient,
		store:           store,
	}
}

// SetupServer создаёт OpenAPI сервер на основе обработчика.
func SetupServer(h *OrderHandler) (*orderv1.Server, error) {
	return orderv1.NewServer(h)
}

// GetOrder реализует операцию getOrder (пример реализации).
// GET /api/v1/orders/{order_uuid}.
func (h *OrderHandler) GetOrder(_ context.Context, params orderv1.GetOrderParams) (orderv1.GetOrderRes, error) {
	// 1. Найти заказ в store (с блокировкой для thread-safety)
	h.store.mu.RLock()
	order, ok := h.store.orders[params.OrderUUID]
	h.store.mu.RUnlock()

	// 2. Если не найден — вернуть 404
	if !ok {
		return &orderv1.GetOrderNotFound{
			Code:    http.StatusNotFound,
			Message: "заказ не найден",
		}, nil
	}

	// 3. Преобразовать в DTO и вернуть
	var shieldUUID orderv1.OptNilUUID
	if order.ShieldUUID != nil {
		shieldUUID = orderv1.NewOptNilUUID(*order.ShieldUUID)
	}

	var weaponUUID orderv1.OptNilUUID
	if order.WeaponUUID != nil {
		weaponUUID = orderv1.NewOptNilUUID(*order.WeaponUUID)
	}

	var transactionUUID orderv1.OptNilUUID
	if order.TransactionUUID != nil {
		transactionUUID = orderv1.NewOptNilUUID(*order.TransactionUUID)
	}

	var paymentMethod orderv1.OptNilPaymentMethod
	if order.PaymentMethod != nil {
		paymentMethod = orderv1.NewOptNilPaymentMethod(orderv1.PaymentMethod(*order.PaymentMethod))
	}

	return &orderv1.OrderDto{
		OrderUUID:       order.OrderUUID,
		HullUUID:        order.HullUUID,
		EngineUUID:      order.EngineUUID,
		ShieldUUID:      shieldUUID,
		WeaponUUID:      weaponUUID,
		TotalPrice:      order.TotalPrice,
		TransactionUUID: transactionUUID,
		PaymentMethod:   paymentMethod,
		Status:          orderv1.OrderStatus(order.Status),
		CreatedAt:       order.CreatedAt,
	}, nil
}

// CreateOrder реализует операцию createOrder
// POST /api/v1/orders.
func (h *OrderHandler) CreateOrder(ctx context.Context, req *orderv1.CreateOrderRequest) (orderv1.CreateOrderRes, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	if req.EngineUUID == uuid.Nil {
		return &orderv1.CreateOrderBadRequest{
			Code:    http.StatusBadRequest,
			Message: "engine_uuid обязателен",
		}, nil
	}

	if req.HullUUID == uuid.Nil {
		return &orderv1.CreateOrderBadRequest{
			Code:    http.StatusBadRequest,
			Message: "hull_uuid обязателен",
		}, nil
	}

	uuids := []string{
		req.EngineUUID.String(),
		req.HullUUID.String(),
	}

	var shieldUUIDPtr *uuid.UUID
	if shieldUUID, ok := req.ShieldUUID.Get(); ok {
		uuids = append(uuids, shieldUUID.String())
		shieldUUIDPtr = &shieldUUID
	}

	var weaponUUIDPtr *uuid.UUID
	if weaponUUID, ok := req.WeaponUUID.Get(); ok {
		uuids = append(uuids, weaponUUID.String())
		weaponUUIDPtr = &weaponUUID
	}

	resp, err := h.inventoryClient.ListParts(ctx, &inventoryv1.ListPartsRequest{
		Uuids: uuids,
	})
	if err != nil {
		st, ok := status.FromError(err)
		if ok {
			switch st.Code() {
			case codes.InvalidArgument:
				return &orderv1.CreateOrderBadRequest{
					Code:    http.StatusBadRequest,
					Message: "некорректный UUID детали",
				}, nil
			case codes.NotFound:
				return &orderv1.CreateOrderNotFound{
					Code:    http.StatusNotFound,
					Message: "деталь не найдена",
				}, nil
			case codes.DeadlineExceeded:
				return &orderv1.CreateOrderInternalServerError{
					Code:    http.StatusInternalServerError,
					Message: "таймаут при обращении к inventory",
				}, nil
			default:
				return &orderv1.CreateOrderInternalServerError{
					Code:    http.StatusInternalServerError,
					Message: "внутренняя ошибка сервера",
				}, nil
			}
		}
		return &orderv1.CreateOrderInternalServerError{
			Code:    http.StatusInternalServerError,
			Message: "внутренняя ошибка сервера",
		}, nil
	}

	if len(resp.Parts) != len(uuids) {
		return &orderv1.CreateOrderNotFound{
			Code:    http.StatusNotFound,
			Message: "деталь не найдена",
		}, nil
	}

	var totalPrice int64
	for _, part := range resp.Parts {
		if part.StockQuantity <= 0 {
			return &orderv1.CreateOrderConflict{
				Code:    http.StatusConflict,
				Message: "деталь отсутствует на складе",
			}, nil
		}

		totalPrice += part.Price
	}

	order := Order{
		OrderUUID:  uuid.New(),
		HullUUID:   req.HullUUID,
		EngineUUID: req.EngineUUID,
		ShieldUUID: shieldUUIDPtr,
		WeaponUUID: weaponUUIDPtr,
		TotalPrice: totalPrice,
		Status:     "PENDING_PAYMENT",
		CreatedAt:  time.Now(),
	}

	h.store.mu.Lock()
	h.store.orders[order.OrderUUID] = order
	h.store.mu.Unlock()

	return &orderv1.CreateOrderResponse{
		OrderUUID:  order.OrderUUID,
		TotalPrice: order.TotalPrice,
	}, nil
}

// PayOrder реализует операцию payOrder
// POST /api/v1/orders/{order_uuid}/pay.
func (h *OrderHandler) PayOrder(ctx context.Context, req *orderv1.PayOrderRequest, params orderv1.PayOrderParams) (orderv1.PayOrderRes, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	h.store.mu.RLock()
	order, ok := h.store.orders[params.OrderUUID]
	h.store.mu.RUnlock()

	if !ok {
		return &orderv1.PayOrderNotFound{
			Code:    http.StatusNotFound,
			Message: "заказ не найден",
		}, nil
	}

	if order.Status == "PAID" {
		return &orderv1.PayOrderConflict{
			Code:    http.StatusConflict,
			Message: "заказ уже оплачен",
		}, nil
	}

	if order.Status == "CANCELLED" {
		return &orderv1.PayOrderConflict{
			Code:    http.StatusConflict,
			Message: "заказ был отменен",
		}, nil
	}

	grpcPaymentMethod, err := mapPaymentMethod(req.PaymentMethod)
	if err != nil {
		return &orderv1.PayOrderBadRequest{
			Code:    http.StatusBadRequest,
			Message: "невалидный payment_method",
		}, nil
	}

	res, err := h.paymentClient.PayOrder(ctx, &paymentv1.PayOrderRequest{
		OrderUuid:     order.OrderUUID.String(),
		PaymentMethod: grpcPaymentMethod,
	})
	if err != nil {
		st, ok := status.FromError(err)
		if ok {
			switch st.Code() {
			case codes.InvalidArgument:
				return &orderv1.PayOrderBadRequest{
					Code:    http.StatusBadRequest,
					Message: "невалидные данные платежа",
				}, nil
			case codes.DeadlineExceeded:
				return &orderv1.PayOrderInternalServerError{
					Code:    http.StatusInternalServerError,
					Message: "таймаут при обращении к payment",
				}, nil
			default:
				return &orderv1.PayOrderInternalServerError{
					Code:    http.StatusInternalServerError,
					Message: "внутренняя ошибка сервера",
				}, nil
			}
		}

		return &orderv1.PayOrderInternalServerError{
			Code:    http.StatusInternalServerError,
			Message: "внутренняя ошибка сервера",
		}, nil
	}

	transactionUuid, err := uuid.Parse(res.TransactionUuid)
	if err != nil {
		return &orderv1.PayOrderInternalServerError{
			Code:    http.StatusInternalServerError,
			Message: "внутренняя ошибка сервера",
		}, nil
	}

	order.Status = "PAID"
	paymentMethod := string(req.PaymentMethod)

	order.PaymentMethod = &paymentMethod
	order.TransactionUUID = &transactionUuid

	h.store.mu.Lock()
	h.store.orders[params.OrderUUID] = order
	h.store.mu.Unlock()

	return &orderv1.PayOrderResponse{
		TransactionUUID: transactionUuid,
	}, nil
}

// CancelOrder реализует операцию cancelOrder
// POST /api/v1/orders/{order_uuid}/cancel.
func (h *OrderHandler) CancelOrder(_ context.Context, params orderv1.CancelOrderParams) (orderv1.CancelOrderRes, error) {
	h.store.mu.RLock()
	order, ok := h.store.orders[params.OrderUUID]
	h.store.mu.RUnlock()

	if !ok {
		return &orderv1.CancelOrderNotFound{
			Code:    http.StatusNotFound,
			Message: "заказ не найден",
		}, nil
	}

	if order.Status == "PAID" {
		return &orderv1.CancelOrderConflict{
			Code:    http.StatusConflict,
			Message: "заказ уже оплачен",
		}, nil
	}

	if order.Status == "CANCELLED" {
		return &orderv1.CancelOrderConflict{
			Code:    http.StatusConflict,
			Message: "заказ был отменен",
		}, nil
	}

	order.Status = "CANCELLED"

	h.store.mu.Lock()
	h.store.orders[params.OrderUUID] = order
	h.store.mu.Unlock()

	return &orderv1.CancelOrderResponse{}, nil
}
