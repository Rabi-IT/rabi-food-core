-- Tenants
INSERT INTO iam.tenants (id, name, slug, language) VALUES
  ('00000000-0000-0000-0000-000000000001', 'Pizzaria Teste', 'pizzaria', 'pt-BR'),
  ('00000000-0000-0000-0000-000000000002', 'Burguer House',  'burguer',  'pt-BR'),
  ('00000000-0000-0000-0000-000000000003', 'Fruit Fresh',   'fruit',    'pt-BR')
ON CONFLICT (id) DO NOTHING;

-- Categories
INSERT INTO catalog.categories (id, tenant_id, name, description) VALUES
  -- pizzaria
  ('10000000-0000-0000-0000-000000000001', '00000000-0000-0000-0000-000000000001', 'Pizzas',     'Pizzas tradicionais e especiais'),
  ('10000000-0000-0000-0000-000000000002', '00000000-0000-0000-0000-000000000001', 'Bebidas',    'Refrigerantes, sucos e cervejas'),
  ('10000000-0000-0000-0000-000000000003', '00000000-0000-0000-0000-000000000001', 'Sobremesas', 'Doces e sobremesas'),
  -- burguer
  ('10000000-0000-0000-0000-000000000004', '00000000-0000-0000-0000-000000000002', 'Burguers',   'Hambúrgueres artesanais'),
  ('10000000-0000-0000-0000-000000000005', '00000000-0000-0000-0000-000000000002', 'Acompanhamentos', 'Fritas, onion rings e mais'),
  ('10000000-0000-0000-0000-000000000006', '00000000-0000-0000-0000-000000000002', 'Bebidas',    'Milkshakes e refrigerantes'),
  -- fruit
  ('10000000-0000-0000-0000-000000000007', '00000000-0000-0000-0000-000000000003', 'Frutas',     'Frutas frescas da estação'),
  ('10000000-0000-0000-0000-000000000008', '00000000-0000-0000-0000-000000000003', 'Sucos',      'Sucos naturais e vitaminas'),
  ('10000000-0000-0000-0000-000000000009', '00000000-0000-0000-0000-000000000003', 'Saladas',    'Saladas de frutas e bowls')
ON CONFLICT (id) DO NOTHING;

-- Products
INSERT INTO catalog.products (id, tenant_id, category_id, name, description, photo, unit, price, is_active, discount_rules) VALUES
  -- pizzaria
  ('20000000-0000-0000-0000-000000000001', '00000000-0000-0000-0000-000000000001', '10000000-0000-0000-0000-000000000001',
   'Pizza Margherita', 'Molho de tomate, mussarela e manjericão fresco', '', 'unidade', 4500, true, '[]'),
  ('20000000-0000-0000-0000-000000000002', '00000000-0000-0000-0000-000000000001', '10000000-0000-0000-0000-000000000001',
   'Pizza Calabresa', 'Molho de tomate, calabresa e cebola', '', 'unidade', 4200, true, '[]'),
  ('20000000-0000-0000-0000-000000000003', '00000000-0000-0000-0000-000000000001', '10000000-0000-0000-0000-000000000001',
   'Pizza 4 Queijos', 'Mussarela, provolone, parmesão e gorgonzola', '', 'unidade', 5000, true, '[]'),
  ('20000000-0000-0000-0000-000000000004', '00000000-0000-0000-0000-000000000001', '10000000-0000-0000-0000-000000000002',
   'Coca-Cola 2L', 'Refrigerante gelado', '', 'unidade', 1200, true, '[]'),
  ('20000000-0000-0000-0000-000000000005', '00000000-0000-0000-0000-000000000001', '10000000-0000-0000-0000-000000000002',
   'Suco de Laranja', 'Suco natural 500ml', '', 'unidade', 900, true, '[]'),
  ('20000000-0000-0000-0000-000000000006', '00000000-0000-0000-0000-000000000001', '10000000-0000-0000-0000-000000000003',
   'Pudim', 'Pudim de leite condensado', '', 'unidade', 800, true, '[]'),
  -- burguer
  ('20000000-0000-0000-0000-000000000007', '00000000-0000-0000-0000-000000000002', '10000000-0000-0000-0000-000000000004',
   'Classic Burguer', 'Blend 180g, queijo cheddar, alface, tomate e maionese', '', 'unidade', 2800, true, '[]'),
  ('20000000-0000-0000-0000-000000000008', '00000000-0000-0000-0000-000000000002', '10000000-0000-0000-0000-000000000004',
   'Bacon Burguer', 'Blend 180g, bacon crocante, queijo prato e molho barbecue', '', 'unidade', 3200, true, '[]'),
  ('20000000-0000-0000-0000-000000000009', '00000000-0000-0000-0000-000000000002', '10000000-0000-0000-0000-000000000004',
   'Double Smash', 'Dois blends 90g smashados, queijo americano e molho especial', '', 'unidade', 3600, true, '[]'),
  ('20000000-0000-0000-0000-000000000010', '00000000-0000-0000-0000-000000000002', '10000000-0000-0000-0000-000000000005',
   'Batata Frita', 'Batatas crocantes temperadas', '', 'porção', 1400, true, '[]'),
  ('20000000-0000-0000-0000-000000000011', '00000000-0000-0000-0000-000000000002', '10000000-0000-0000-0000-000000000005',
   'Onion Rings', 'Anéis de cebola empanados', '', 'porção', 1600, true, '[]'),
  ('20000000-0000-0000-0000-000000000012', '00000000-0000-0000-0000-000000000002', '10000000-0000-0000-0000-000000000006',
   'Milkshake Chocolate', 'Milkshake cremoso de chocolate 400ml', '', 'unidade', 1800, true, '[]'),
  -- fruit
  ('20000000-0000-0000-0000-000000000013', '00000000-0000-0000-0000-000000000003', '10000000-0000-0000-0000-000000000007',
   'Morango', 'Morango fresco da estação', '', 'kg', 1200, true, '[]'),
  ('20000000-0000-0000-0000-000000000014', '00000000-0000-0000-0000-000000000003', '10000000-0000-0000-0000-000000000007',
   'Manga Palmer', 'Manga doce e suculenta', '', 'kg', 900, true, '[]'),
  ('20000000-0000-0000-0000-000000000015', '00000000-0000-0000-0000-000000000003', '10000000-0000-0000-0000-000000000007',
   'Uva Itália', 'Uva verde sem semente', '', 'kg', 1500, true, '[]'),
  ('20000000-0000-0000-0000-000000000016', '00000000-0000-0000-0000-000000000003', '10000000-0000-0000-0000-000000000008',
   'Suco de Acerola', 'Suco natural de acerola 500ml', '', 'unidade', 700, true, '[]'),
  ('20000000-0000-0000-0000-000000000017', '00000000-0000-0000-0000-000000000003', '10000000-0000-0000-0000-000000000008',
   'Vitamina de Banana', 'Vitamina com banana, leite e mel 400ml', '', 'unidade', 900, true, '[]'),
  ('20000000-0000-0000-0000-000000000018', '00000000-0000-0000-0000-000000000003', '10000000-0000-0000-0000-000000000009',
   'Salada de Frutas', 'Mix de frutas frescas com granola', '', 'unidade', 1100, true, '[]')
ON CONFLICT (id) DO NOTHING;
