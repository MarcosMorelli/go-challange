db = db.getSiblingDB("jobsity");

db.createCollection("users");

db.users.insertOne({
  email: "user1@jobsity.com",
  password: "password",
  created_at: new Date(),
});

db.users.insertOne({
  email: "user2@jobsity.com",
  password: "password",
  created_at: new Date(),
});

db.channels.insertOne({
  name: "channel-1",
  description: "Channel 1",
  created_by: "user1@jobsity.com",
  created_at: new Date(),
});

db.channels.insertOne({
  name: "channel-2",
  description: "Channel 2",
  created_by: "user1@jobsity.com",
  created_at: new Date(),
});
