#!/usr/bin/node

const secret = "secret";

const crypto = require("crypto");
const fs = require("fs");
const readline = require("readline");

const rl = readline.createInterface({
  input: process.stdin,
  output: process.stdout,
});

let username = "";
let password = "";
let roles = [];

rl.question("Choose username: ", (u) => {
  username = u;
  rl.question("Choose password: ", (p) => {
    password = p;
    rl.question("Assign roles: ", (r) => {
      roles = r.split(",");
      rl.close();

      let users = {};

      if (fs.existsSync("./data/users.json")) {
        users = JSON.parse(fs.readFileSync("./data/users.json", "utf8"));
      }

      const hashedPassword = crypto
        .createHmac("sha256", secret)
        .update(password)
        .digest("hex");

      users[username] = {
        username,
        password: hashedPassword,
        roles,
      };

      fs.writeFileSync("./data/users.json", JSON.stringify(users, null, 2));

      console.log(`User ${username} has been added!`);
    });
  });
});
