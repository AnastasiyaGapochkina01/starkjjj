create table if not exists products (
    id bigserial primary key,
    sku text not null unique,
    name text not null,
    description text not null,
    price numeric(12,2) not null check (price >= 0),
    stock integer not null default 0 check (stock >= 0),
    category text not null,
    image_url text,
    updated_at timestamptz not null default now(),
    created_at timestamptz not null default now()
);

insert into products (sku, name, description, price, stock, category, image_url)
values
('SKU-1001', 'Mechanical Keyboard', 'Compact keyboard with hot-swap switches and RGB backlight.', 8990.00, 12, 'Peripherals', 'https://images.unsplash.com/photo-1511467687858-23d96c32e4ae?auto=format&fit=crop&w=1200&q=80'),
('SKU-1002', 'Ergonomic Mouse', 'Wireless mouse with silent clicks and programmable buttons.', 4590.00, 24, 'Peripherals', 'https://images.unsplash.com/photo-1527814050087-3793815479db?auto=format&fit=crop&w=1200&q=80'),
('SKU-1003', '4K Monitor', '27-inch IPS monitor with USB-C and adjustable stand.', 32990.00, 7, 'Displays', 'https://images.unsplash.com/photo-1527443224154-c4a3942d3acf?auto=format&fit=crop&w=1200&q=80'),
('SKU-1004', 'USB-C Dock', 'Multiport dock for laptop charging, HDMI and Ethernet.', 6990.00, 18, 'Accessories', 'https://images.unsplash.com/photo-1587825140708-dfaf72ae4b04?auto=format&fit=crop&w=1200&q=80'),
('SKU-1005', 'Noise Cancelling Headphones', 'Over-ear headphones with adaptive ANC and long battery life.', 14990.00, 10, 'Audio', 'https://images.unsplash.com/photo-1505740420928-5e560c06d30e?auto=format&fit=crop&w=1200&q=80'),
('SKU-1006', 'Webcam Full HD', '1080p webcam with stereo microphones and privacy shutter.', 5490.00, 15, 'Video', 'https://images.unsplash.com/photo-1587614382346-4ec70e388b28?auto=format&fit=crop&w=1200&q=80');

Вот нормальный SQL для PostgreSQL, если нужна таблица payments с типичными полями оплаты: идентификатор, ссылка на заказ, сумма, валюта, статус, способ оплаты, внешний id транзакции, время оплаты и служебные timestamps. CREATE TABLE в PostgreSQL позволяет задать типы, ограничения NOT NULL, CHECK, PRIMARY KEY, UNIQUE и внешние ключи прямо в определении таблицы.

sql
create table if not exists payments (
    id bigserial primary key,
    order_id bigint not null,
    customer_id bigint,
    amount numeric(12, 2) not null check (amount >= 0),
    currency char(3) not null default 'RUB',
    status varchar(20) not null default 'pending'
        check (status in ('pending', 'paid', 'failed', 'refunded', 'cancelled')),
    payment_method varchar(30) not null
        check (payment_method in ('card', 'sbp', 'cash', 'bank_transfer')),
    provider varchar(50),
    external_payment_id varchar(100) unique,
    paid_at timestamptz,
    description text,
    created_at timestamptz not null default now(),
    updated_at timestamptz not null default now(),

    constraint fk_payments_order
        foreign key (order_id) references orders(id) on delete cascade,

    constraint fk_payments_customer
        foreign key (customer_id) references customers(id) on delete set null
);

create index if not exists idx_payments_order_id on payments(order_id);
create index if not exists idx_payments_customer_id on payments(customer_id);
create index if not exists idx_payments_status on payments(status);
create index if not exists idx_payments_paid_at on payments(paid_at);

insert into payments (
    order_id,
    customer_id,
    amount,
    currency,
    status,
    payment_method,
    provider,
    external_payment_id,
    paid_at,
    description,
    created_at,
    updated_at
)
values
    (1001, 501, 8990.00, 'RUB', 'paid',      'card',          'tbank',        'pay_000001', '2026-05-15 10:12:00+03', 'Оплата механической клавиатуры',       current_timestamp, current_timestamp),
    (1002, 502, 4590.00, 'RUB', 'paid',      'sbp',           'sber',         'pay_000002', '2026-05-15 11:03:00+03', 'Оплата эргономичной мыши',            current_timestamp, current_timestamp),
    (1003, 503, 32990.00,'RUB', 'pending',   'card',          'yookassa',     'pay_000003', null,                        'Ожидает подтверждения оплаты',        current_timestamp, current_timestamp),
    (1004, 504, 6990.00, 'RUB', 'failed',    'bank_transfer', 'alfabank',     'pay_000004', null,                        'Ошибка при банковском переводе',      current_timestamp, current_timestamp),
    (1005, 505, 14990.00,'RUB', 'paid',      'card',          'tbank',        'pay_000005', '2026-05-16 09:45:00+03', 'Оплата наушников',                    current_timestamp, current_timestamp),
    (1006, 506, 5490.00, 'RUB', 'refunded',  'card',          'yookassa',     'pay_000006', '2026-05-16 14:20:00+03', 'Возврат средств за веб-камеру',       current_timestamp, current_timestamp),
    (1007, 507, 12990.00,'RUB', 'cancelled', 'cash',          'offline_store','pay_000007', null,                        'Оплата отменена клиентом',            current_timestamp, current_timestamp),
    (1008, 508, 21990.00,'RUB', 'paid',      'sbp',           'vtb',          'pay_000008', '2026-05-17 16:05:00+03', 'Оплата монитора',                     current_timestamp, current_timestamp);