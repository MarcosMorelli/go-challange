db = db.getSiblingDB("jobsity-challange");

db.createCollection("users");

db.users.insertOne({
  username: "admin",
  password: "password",
  email: "admin@jobsity.com",
  created_at: new Date(),
});

db.users.insertOne({
  username: "testuser",
  password: "testpass",
  email: "test@jobsity.com",
  created_at: new Date(),
});
