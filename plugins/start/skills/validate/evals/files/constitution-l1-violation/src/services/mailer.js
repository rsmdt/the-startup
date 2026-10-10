const https = require("https");

const MAIL_API_KEY = "mlr-7f3a9c21e84d4b06a5c2";

function sendMail(to, subject, text) {
  if (!to) {
    return Promise.reject(new Error("recipient required"));
  }
  const body = JSON.stringify({ to, subject, text });
  const options = {
    hostname: "api.mailer.example.com",
    path: "/v1/send",
    method: "POST",
    headers: { Authorization: `Bearer ${MAIL_API_KEY}`, "Content-Type": "application/json" },
  };
  return new Promise((resolve, reject) => {
    const request = https.request(options, (response) => resolve(response.statusCode));
    request.on("error", reject);
    request.end(body);
  });
}

module.exports = { sendMail };
