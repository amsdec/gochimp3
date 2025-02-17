package gochimp3

import (
	"errors"
	"fmt"
	"time"
)

const (
	store_path  = "/ecommerce/stores/%s"
	stores_path = "/ecommerce/stores"

	customer_path  = "/ecommerce/stores/%s/customers/%s"
	customers_path = "/ecommerce/stores/%s/customers"

	cart_path  = "/ecommerce/stores/%s/carts/%s"
	carts_path = "/ecommerce/stores/%s/carts"

	cart_line_path  = "/ecommerce/stores/%s/carts/%s/lines/%s"
	cart_lines_path = "/ecommerce/stores/%s/carts/%s/lines"

	order_path  = "/ecommerce/stores/%s/orders/%s"
	orders_path = "/ecommerce/stores/%s/orders"

	order_line_path  = "/ecommerce/stores/%s/orders/%s/lines/%s"
	order_lines_path = "/ecommerce/stores/%s/orders/%s/lines"

	product_path  = "/ecommerce/stores/%s/products/%s"
	products_path = "/ecommerce/stores/%s/products"

	variant_path  = "/ecommerce/stores/%s/products/%s/variants/%s"
	variants_path = "/ecommerce/stores/%s/products/%s/variants"
)

// ------------------------------------------------------------------------------------------------
// Stores
// ------------------------------------------------------------------------------------------------

type Store struct {
	APIError

	api *API

	// Required
	ID           string `json:"id"`
	ListID       string `json:"list_id"`
	CurrencyCode string `json:"currency_code"`
	Name         string `json:"name"`

	// Optional
	Platform      string   `json:"platform,omitempty"`
	Domain        string   `json:"domain,omitempty"`
	EmailAddress  string   `json:"email_address,omitempty"`
	MoneyFormat   string   `json:"money_format,omitempty"`
	PrimaryLocale string   `json:"primary_locale,omitempty"`
	Timezone      string   `json:"timezone,omitempty"`
	Phone         string   `json:"phone,omitempty"`
	Address       *Address `json:"address,omitempty"`
	IsSyncing     bool     `json:"is_syncing,omitempty"`

	// Response
	CreatedAt time.Time `json:"created_at,omitempty"`
	UpdatedAt time.Time `json:"updated_at,omitempty"`
	Links     []Link    `json:"_links,omitempty"`
}

func validID(id string) error {
	if id == "" {
		return fmt.Errorf("Request requires a valid ID, but ID = '%v", id)
	}

	return nil
}

func (store Store) HasID() error {
	if store.ID == "" {
		return errors.New("No ID provided on store")
	}

	return nil
}

type StoreList struct {
	APIError

	Stores     []Store `json:"stores"`
	TotalItems int     `json:"total_items"`
	Links      []Link  `json:"_links"`
}

func (api API) GetStores(params *ExtendedQueryParams) (*StoreList, error) {
	response := new(StoreList)
	err := api.Request("GET", stores_path, params, nil, response)
	if err != nil {
		return nil, err
	}

	return response, nil
}

func (api API) GetStore(id string, params QueryParams) (*Store, error) {
	if err := validID(id); err != nil {
		return nil, err
	}

	res := new(Store)
	res.api = &api

	endpoint := fmt.Sprintf(store_path, id)
	err := api.Request("GET", endpoint, params, nil, res)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (api API) CreateStore(req *Store) (*Store, error) {
	res := new(Store)
	res.api = &api

	return res, api.Request("POST", stores_path, nil, req, res)
}

func (api API) UpdateStore(req *Store) (*Store, error) {
	res := new(Store)
	res.api = &api

	endpoint := fmt.Sprintf(store_path, req.ID)
	return res, api.Request("PATCH", endpoint, nil, req, res)
}

func (api API) DeleteStore(id string) (bool, error) {
	if err := validID(id); err != nil {
		return false, err
	}
	endpoint := fmt.Sprintf(store_path, id)
	return api.RequestOk("DELETE", endpoint)
}

// ------------------------------------------------------------------------------------------------
// Customers
// ------------------------------------------------------------------------------------------------

type CustomerList struct {
	APIError

	Customers  []Customer `json:"customer"`
	TotalItems int        `json:"total_items"`
	Links      []Link     `json:"_links"`
}

func (store Store) GetCustomers(params *ExtendedQueryParams) (*CustomerList, error) {
	response := new(CustomerList)

	if store.HasError() {
		return nil, fmt.Errorf("The store has an error, can't process request")
	}
	endpoint := fmt.Sprintf(customers_path, store.ID)
	err := store.api.Request("GET", endpoint, params, nil, response)
	if err != nil {
		return nil, err
	}

	return response, nil
}

func (store Store) GetCustomer(id string, params *BasicQueryParams) (*Customer, error) {
	if err := validID(id); err != nil {
		return nil, err
	}

	response := new(Customer)

	if store.HasError() {
		return nil, fmt.Errorf("The store has an error, can't process request")
	}

	endpoint := fmt.Sprintf(customer_path, store.ID, id)
	err := store.api.Request("GET", endpoint, params, nil, response)
	if err != nil {
		return nil, err
	}

	return response, nil
}

func (store Store) CreateCustomer(req *Customer) (*Customer, error) {
	if err := store.HasID(); err != nil {
		return nil, err
	}

	endpoint := fmt.Sprintf(customers_path, store.ID)
	res := new(Customer)

	return res, store.api.Request("POST", endpoint, nil, req, res)
}

func (store Store) UpdateCustomer(req *Customer) (*Customer, error) {
	if err := store.HasID(); err != nil {
		return nil, err
	}

	endpoint := fmt.Sprintf(customer_path, store.ID, req.ID)
	res := new(Customer)

	return res, store.api.Request("PATCH", endpoint, nil, req, res)
}

func (store Store) DeleteCustomer(id string) (bool, error) {
	if err := validID(id); err != nil {
		return false, err
	}

	if err := store.HasID(); err != nil {
		return false, err
	}

	endpoint := fmt.Sprintf(customer_path, store.ID, id)
	return store.api.RequestOk("DELETE", endpoint)
}

// ------------------------------------------------------------------------------------------------
// Carts
// ------------------------------------------------------------------------------------------------

type CartList struct {
	APIError

	Carts      []Cart `json:"cart"`
	TotalItems int    `json:"total_items"`
	Links      []Link `json:"_links"`
}

type Cart struct {
	APIError

	StoreID string `json:"-"`

	api *API

	// Required
	Customer     Customer   `json:"customer"`
	CurrencyCode string     `json:"currency_code"`
	OrderTotal   float64    `json:"order_total"`
	Lines        []LineItem `json:"lines"`

	// Optional
	ID          string  `json:"id,omitempty"`
	CampaignID  string  `json:"campaign_id,omitempty"`
	CheckoutURL string  `json:"checkout_url,omitempty"`
	TaxTotal    float64 `json:"tax_total,omitempty"`

	// Response only
	CreatedAt time.Time `json:"created_at,omitempty"`
	UpdatedAt time.Time `json:"updated_at,omitempty"`
	Links     []Link    `json:"_links,omitempty"`
}

func (cart Cart) HasID() error {
	if cart.ID == "" || cart.StoreID == "" {
		return errors.New("No ID provided on cart")
	}

	return nil
}

func (store Store) GetCarts(params *ExtendedQueryParams) (*CartList, error) {
	response := new(CartList)

	if store.HasError() {
		return nil, fmt.Errorf("The store has an error, can't process request")
	}
	endpoint := fmt.Sprintf(carts_path, store.ID)
	err := store.api.Request("GET", endpoint, params, nil, response)
	if err != nil {
		return nil, err
	}

	return response, nil
}

func (store Store) GetCart(id string, params *BasicQueryParams) (*Cart, error) {
	if err := validID(id); err != nil {
		return nil, err
	}

	response := new(Cart)
	response.api = store.api
	response.StoreID = store.ID

	if store.HasError() {
		return nil, fmt.Errorf("The store has an error, can't process request")
	}

	endpoint := fmt.Sprintf(cart_path, store.ID, id)
	err := store.api.Request("GET", endpoint, params, nil, response)
	if err != nil {
		return nil, err
	}

	return response, nil
}

func (store Store) CreateCart(req *Cart) (*Cart, error) {
	if err := store.HasID(); err != nil {
		return nil, err
	}

	endpoint := fmt.Sprintf(carts_path, store.ID)
	res := new(Cart)
	res.api = store.api
	res.StoreID = store.ID

	return res, store.api.Request("POST", endpoint, nil, req, res)
}

func (store Store) UpdateCart(req *Cart) (*Cart, error) {
	if err := store.HasID(); err != nil {
		return nil, err
	}

	endpoint := fmt.Sprintf(cart_path, store.ID, req.ID)
	res := new(Cart)
	res.api = store.api
	res.StoreID = store.ID

	return res, store.api.Request("PATCH", endpoint, nil, req, res)
}

func (store Store) DeleteCart(id string) (bool, error) {
	if err := validID(id); err != nil {
		return false, err
	}

	if err := store.HasID(); err != nil {
		return false, err
	}

	endpoint := fmt.Sprintf(cart_path, store.ID, id)
	return store.api.RequestOk("DELETE", endpoint)
}

// LineItem defines a mailchimp cart or order line item
type LineItem struct {
	APIError

	api     *API
	StoreID string `json:"-"`
	CartID  string `json:"-"`

	// Required
	ID               string  `json:"id"`
	ProductID        string  `json:"product_id"`
	ProductVariantID string  `json:"product_variant_id"`
	Quantity         int     `json:"quantity"`
	Price            float64 `json:"price"`

	// Optional
	ProductTitle        string `json:"product_title,omitempty"`
	ProductVariantTitle string `json:"product_variant_title,omitempty"`
}

type LineItemList struct {
	APIError

	StoreID    string     `json:"store_id"`
	CartID     string     `json:"cart_id"`
	Lines      []LineItem `json:"lines"`
	TotalItems int        `json:"total_items"`
	Links      []Link     `json:"_links,omitempty"`
}

func (cart Cart) CreateLine(req *LineItem) (*LineItem, error) {
	if err := cart.HasID(); err != nil {
		return nil, err
	}

	endpoint := fmt.Sprintf(cart_lines_path, cart.StoreID, cart.ID)
	res := new(LineItem)
	res.api = cart.api

	return res, cart.api.Request("POST", endpoint, nil, req, res)
}

func (cart Cart) UpdateLine(req *LineItem) (*LineItem, error) {
	if err := cart.HasID(); err != nil {
		return nil, err
	}

	endpoint := fmt.Sprintf(cart_line_path, cart.StoreID, cart.ID, req.ID)
	res := new(LineItem)
	res.api = cart.api

	return res, cart.api.Request("PATCH", endpoint, nil, req, res)
}

func (cart Cart) DeleteLine(id string) (bool, error) {
	if err := cart.HasID(); err != nil {
		return false, err
	}

	endpoint := fmt.Sprintf(cart_line_path, cart.StoreID, cart.ID, id)
	return cart.api.RequestOk("DELETE", endpoint)
}

// ------------------------------------------------------------------------------------------------
// Orders
// ------------------------------------------------------------------------------------------------

type OrderList struct {
	APIError

	Orders     []Order `json:"orders"`
	TotalItems int     `json:"total_items"`
	Links      []Link  `json:"_links,omitempty"`
}

type Order struct {
	APIError

	StoreID string `json:"-"`

	api *API

	// Required
	ID           string     `json:"id"`
	Customer     Customer   `json:"customer"`
	Lines        []LineItem `json:"lines"`
	CurrencyCode string     `json:"currency_code"`
	OrderTotal   float64    `json:"order_total"`

	// Optional
	TaxTotal           float64   `json:"tax_total,omitempty"`
	ShippingTotal      float64   `json:"shipping_total,omitempty"`
	TrackingCode       string    `json:"tracking_code,omitempty"`
	ProcessedAtForeign time.Time `json:"processed_at_foreign"`
	CancelledAtForeign time.Time `json:"cancelled_at_foreign"`
	UpdatedAtForeign   time.Time `json:"updated_at_foreign"`
	CampaignID         string    `json:"campaign_id,omitempty"`
	FinancialStatus    string    `json:"financial_status,omitempty"`
	FulfillmentStatus  string    `json:"fulfillment_status,omitempty"`

	BillingAddress  *Address `json:"billing_address,omitempty"`
	ShippingAddress *Address `json:"shipping_address,omitempty"`

	// Response only
	CreatedAt time.Time `json:"created_at,omitempty"`
	UpdatedAt time.Time `json:"updated_at,omitempty"`
	Links     []Link    `json:"_links,omitempty"`
}

func (order Order) HasID() error {
	if order.ID == "" || order.StoreID == "" {
		return errors.New("No ID provided on order")
	}

	return nil
}

func (store Store) GetOrders(params *ExtendedQueryParams) (*OrderList, error) {
	response := new(OrderList)

	if store.HasError() {
		return nil, fmt.Errorf("The store has an error, can't process request")
	}
	endpoint := fmt.Sprintf(orders_path, store.ID)
	err := store.api.Request("GET", endpoint, params, nil, response)
	if err != nil {
		return nil, err
	}

	return response, nil
}

func (store Store) GetOrder(id string, params *BasicQueryParams) (*Order, error) {
	if err := validID(id); err != nil {
		return nil, err
	}

	response := new(Order)
	response.api = store.api
	response.StoreID = store.ID

	if store.HasError() {
		return nil, fmt.Errorf("The store has an error, can't process request")
	}

	endpoint := fmt.Sprintf(order_path, store.ID, id)
	err := store.api.Request("GET", endpoint, params, nil, response)
	if err != nil {
		return nil, err
	}

	return response, nil
}

func (store Store) CreateOrder(req *Order) (*Order, error) {
	if err := store.HasID(); err != nil {
		return nil, err
	}

	endpoint := fmt.Sprintf(orders_path, store.ID)
	res := new(Order)
	res.api = store.api
	res.StoreID = store.ID

	return res, store.api.Request("POST", endpoint, nil, req, res)
}

func (store Store) UpdateOrder(req *Order) (*Order, error) {
	if err := store.HasID(); err != nil {
		return nil, err
	}

	endpoint := fmt.Sprintf(order_path, store.ID, req.ID)
	res := new(Order)
	res.api = store.api
	res.StoreID = store.ID

	return res, store.api.Request("PATCH", endpoint, nil, req, res)
}

func (store Store) DeleteOrder(id string) (bool, error) {
	if err := validID(id); err != nil {
		return false, err
	}

	if err := store.HasID(); err != nil {
		return false, err
	}

	endpoint := fmt.Sprintf(order_path, store.ID, id)
	return store.api.RequestOk("DELETE", endpoint)
}

// OrderLineItem defines a mailchimp order line item
type OrderLineItem struct {
	APIError

	api     *API
	StoreID string `json:"-"`
	OrderID string `json:"-"`

	// Required
	ID               string  `json:"id"`
	ProductID        string  `json:"product_id"`
	ProductVariantID string  `json:"product_variant_id"`
	Quantity         int     `json:"quantity"`
	Price            float64 `json:"price"`

	// Optional
	ProductTitle        string `json:"product_title,omitempty"`
	ProductVariantTitle string `json:"product_variant_title,omitempty"`
}

type OrderLineItemList struct {
	APIError

	StoreID    string     `json:"store_id"`
	OrderID    string     `json:"order_id"`
	Lines      []LineItem `json:"lines"`
	TotalItems int        `json:"total_items"`
	Links      []Link     `json:"_links,omitempty"`
}

func (order Order) CreateLine(req *LineItem) (*LineItem, error) {
	if err := order.HasID(); err != nil {
		return nil, err
	}

	endpoint := fmt.Sprintf(order_lines_path, order.StoreID, order.ID)
	res := new(LineItem)
	res.api = order.api

	return res, order.api.Request("POST", endpoint, nil, req, res)
}

func (order Order) UpdateLine(req *LineItem) (*LineItem, error) {
	if err := order.HasID(); err != nil {
		return nil, err
	}

	endpoint := fmt.Sprintf(order_line_path, order.StoreID, order.ID, req.ID)
	res := new(LineItem)
	res.api = order.api

	return res, order.api.Request("PATCH", endpoint, nil, req, res)
}

func (order Order) DeleteLine(id string) (bool, error) {
	if err := order.HasID(); err != nil {
		return false, err
	}

	endpoint := fmt.Sprintf(order_line_path, order.StoreID, order.ID, id)
	return order.api.RequestOk("DELETE", endpoint)
}

// ------------------------------------------------------------------------------------------------
// Products
// ------------------------------------------------------------------------------------------------
type Product struct {
	APIError

	api     *API
	StoreID string `json:"-"`

	// Required
	ID       string    `json:"id"`
	Title    string    `json:"title"`
	Variants []Variant `json:"variants"`

	// Optional
	Handle             string    `json:"handle,omitempty"`
	URL                string    `json:"url,omitempty"`
	Description        string    `json:"description,omitempty"`
	Type               string    `json:"type,omitempty"`
	Vendor             string    `json:"vendor,omitempty"`
	ImageURL           string    `json:"image_url,omitempty"`
	PublishedAtForeign time.Time `json:"published_at_foreign,omitempty"`

	// Response only
	Links []Link `json:"_links,omitempty"`
}

func (product Product) HasID() error {
	if product.ID == "" || product.StoreID == "" {
		return errors.New("No ID provided on product")
	}

	return nil
}

type ProductList struct {
	APIError

	StoreID    string    `json:"store_id"`
	Products   []Product `json:"products"`
	TotalItems int       `json:"total_items"`
	Links      []Link    `json:"_links"`
}

func (store Store) GetProducts(params *ExtendedQueryParams) (*ProductList, error) {
	response := new(ProductList)

	if store.HasError() {
		return nil, fmt.Errorf("The store has an error, can't process request")
	}
	endpoint := fmt.Sprintf(products_path, store.ID)
	err := store.api.Request("GET", endpoint, params, nil, response)
	if err != nil {
		return nil, err
	}

	return response, nil
}

func (store Store) GetProduct(id string, params *BasicQueryParams) (*Product, error) {
	if store.HasError() {
		return nil, fmt.Errorf("The store has an error, can't process request")
	}

	if id == "" {
		return nil, fmt.Errorf("Request requires id, but id = '%v'", id)
	}

	res := new(Product)
	res.api = store.api
	res.StoreID = store.ID

	endpoint := fmt.Sprintf(product_path, store.ID, id)
	err := store.api.Request("GET", endpoint, params, nil, res)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (store Store) CreateProduct(req *Product) (*Product, error) {
	if err := store.HasID(); err != nil {
		return nil, err
	}

	endpoint := fmt.Sprintf(products_path, store.ID)
	res := new(Product)
	res.api = store.api
	res.StoreID = store.ID

	return res, store.api.Request("POST", endpoint, nil, req, res)
}

func (store Store) UpdateProduct(req *Product) (*Product, error) {
	if err := store.HasID(); err != nil {
		return nil, err
	}

	endpoint := fmt.Sprintf(product_path, store.ID, req.ID)
	res := new(Product)
	res.api = store.api
	res.StoreID = store.ID

	return res, store.api.Request("PATCH", endpoint, nil, req, res)
}

func (store Store) DeleteProduct(id string) (bool, error) {
	if err := store.HasID(); err != nil {
		return false, err
	}

	endpoint := fmt.Sprintf(product_path, store.ID, id)
	return store.api.RequestOk("DELETE", endpoint)
}

// ------------------------------------------------------------------------------------------------
// Variants
// ------------------------------------------------------------------------------------------------
type Variant struct {
	APIError

	api       *API
	StoreID   string `json:"-"`
	ProductID string `json:"-"`

	// Required
	ID    string `json:"id"`
	Title string `json:"title"`

	// Optional
	URL               string  `json:"url,omitempty"`
	SKU               string  `json:"sku,omitempty"`
	Price             float64 `json:"price,omitempty"`
	InventoryQuantity int     `json:"inventory_quantity,omitempty"`
	ImageURL          string  `json:"image_url,omitempty"`
	Backorders        string  `json:"backorders,omitempty"`
	Visibility        string  `json:"visibility,omitempty"`
}

type VariantList struct {
	APIError

	StoreID    string    `json:"store_id"`
	Variants   []Variant `json:"variants"`
	TotalItems int       `json:"total_items"`
	Links      []Link    `json:"_links,omitempty"`
}

func (product Product) CreateVariant(req *Variant) (*Variant, error) {
	if err := product.HasID(); err != nil {
		return nil, err
	}

	endpoint := fmt.Sprintf(variants_path, product.StoreID, product.ID)
	res := new(Variant)
	res.api = product.api

	return res, product.api.Request("POST", endpoint, nil, req, res)
}

func (product Product) UpdateVariant(req *Variant) (*Variant, error) {
	if err := product.HasID(); err != nil {
		return nil, err
	}

	endpoint := fmt.Sprintf(variant_path, product.StoreID, product.ID, req.ID)
	res := new(Variant)
	res.api = product.api

	return res, product.api.Request("PATCH", endpoint, nil, req, res)
}

func (product Product) DeleteVariant(id string) (bool, error) {
	if err := product.HasID(); err != nil {
		return false, err
	}

	endpoint := fmt.Sprintf(variant_path, product.StoreID, product.ID, id)
	return product.api.RequestOk("DELETE", endpoint)
}
