const TOKEN_STORAGE_LOCATION = "poet_token";
let authToken;

async function getPoetInfo() {
  const poet_token = localStorage.getItem(TOKEN_STORAGE_LOCATION);
  const headers = poet_token ? { poet_token: poet_token } : undefined;

  const res = await fetch("/poetry/poet/auth", {
    method: "GET",
    headers: headers,
  });
  if (!res.ok) {
    throw new Error("Not an OK response");
  }
  const resJSON = await res.json();
  localStorage.setItem(TOKEN_STORAGE_LOCATION, resJSON.token);
  authToken = resJSON.token;

  const event = new CustomEvent("tokenObtained", {
    bubbles: true,
    cancelable: true,
  });
  document.body.dispatchEvent(event);
  return;
}

const auth = getPoetInfo();

htmx.on("htmx:confirm", (e) => {
  if (authToken == null) {
    e.preventDefault();
    auth.then(() => e.detail.issueRequest());
  }
});

htmx.on("htmx:configRequest", (e) => {
  e.detail.headers["poet_token"] = authToken;
});
