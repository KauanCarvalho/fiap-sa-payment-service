// Development database.
db = db.getSiblingDB('fiap_sa_payment_service');
db.payments.insertMany([
    { amount: 10 }
]);

// Test database.
testDb = db.getSiblingDB('fiap_sa_payment_service_test');
testDb.products.insertMany([
  { amount: 10 }
]);
