package usecase

type Services struct {
	Auth      *AuthService
	User      *UserService
	Product   *ProductService
	Order     *OrderService
	OrderItem *OrderItemService
}
