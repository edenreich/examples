const fs = require("fs").promises;
const cr = require("crypto");

const secret = "secret";

async function login(r) {
  if (r.method !== "POST") {
    r.return(
      405,
      JSON.stringify({
        status: "INVALID",
        message: "Method not allowed",
      })
    );
    return;
  }

  const body = JSON.parse(r.requestBody);
  const username = body.username;
  const password = body.password;
  const users = await fs
    .readFile("/etc/nginx/data/users.json", "utf8")
    .then((u) => JSON.parse(u));
  if (!users[username]) {
    r.return(
      403,
      JSON.stringify({
        status: "INVALID",
        message: "User not found",
      })
    );
  }

  try {
    const user = users[username];
    if (
      cr.createHmac("sha256", secret).update(password).digest("hex") ===
      user.password
    ) {
      const token = await jwt(user);
      const reply = {
        status: "OK",
        token,
      };
      r.return(200, JSON.stringify(reply));
    } else {
      r.return(
        403,
        JSON.stringify({
          status: "INVALID",
          message: "Invalid username or password",
        })
      );
    }
  } catch (e) {
    r.return(500, e.message);
    return;
  }
}

async function generate_hs256_jwt(init_claims, secret, valid) {
  let header = { typ: "JWT", alg: "HS256" };
  let claims = Object.assign(init_claims, {
    exp: Math.floor(Date.now() / 1000) + valid,
  });

  let data = [header, claims]
    .map(JSON.stringify)
    .map((v) => Buffer.from(v).toString("base64url"))
    .join(".");

  let wc_key = await crypto.subtle.importKey(
    "raw",
    secret,
    { name: "HMAC", hash: "SHA-256" },
    false,
    ["sign"]
  );
  let sign = await crypto.subtle.sign({ name: "HMAC" }, wc_key, data);

  return data + "." + Buffer.from(sign).toString("base64url");
}

function jwt(user) {
  const claims = {
    iss: "nginx",
    sub: user.username,
    roles: user.roles,
  };
  return Promise.resolve(generate_hs256_jwt(claims, secret, 3600));
}

async function verify(r) {
  if (!r.headersIn.Authorization) {
    r.return(
      403,
      JSON.stringify({
        status: "INVALID",
        message: "Authorization header is missing",
      })
    );
    return;
  }
  const token = r.headersIn.Authorization.split(" ")[1];
  const parts = token.split(".");
  const header_b64 = parts[0];
  const payload_b64 = parts[1];
  const signature_b64 = parts[2];

  if (parts.length !== 3) {
    r.return(
      403,
      JSON.stringify({
        status: "INVALID",
        message: "JWT has too many or too little segments",
      })
    );
  }

  const header_json = Buffer.from(header_b64, "base64url").toString();
  const header_obj = JSON.parse(header_json);

  if (header_obj.alg === "HS256") {
    const wc_key = await crypto.subtle.importKey(
      "raw",
      secret,
      { name: "HMAC", hash: "SHA-256" },
      false,
      ["verify"]
    );
    const valid = await crypto.subtle.verify(
      { name: "HMAC" },
      wc_key,
      Buffer.from(signature_b64, "base64url"),
      Buffer.from(header_b64 + "." + payload_b64)
    );
    if (valid) {
      r.internalRedirect("@backend");
    } else {
      r.return(
        403,
        JSON.stringify({ status: "INVALID", message: "Failed to verify JWT" })
      );
    }
  } else {
    r.return(
      403,
      JSON.stringify({ status: "INVALID", message: "Algorithm not supported" })
    );
  }
}

export default { verify, login };
