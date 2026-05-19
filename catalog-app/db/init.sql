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
