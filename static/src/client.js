class ApiClient {
  static default = new ApiClient("/api");

  constructor(baseURL) {
    this.baseURL = baseURL;
  }

  async request(method = "GET", path, data) {
    const headers = {
      "Content-Type": "application/json",
      Accept: "application/json",
    };

    const res = await fetch(`${this.baseURL}${path}`, {
      method,
      headers,
      body: data ? JSON.stringify(data) : null,
    });

    if (res.ok) {
      return res.json();
    }

    throw new ApiError(res.status, res.statusText);
  }

  mailboxes() {
    return this.request("GET", "/mailboxes");
  }

  threads(folder, page = 0, perPage = 50) {
    return this.request(
      "GET",
      `/mailboxes/${folder}?page=${page}&per=${perPage}`
    );
  }

  thread(threadId) {
    return this.request("GET", `/threads/${threadId}`);
  }
}

export class ApiError extends Error {
  constructor(status, statusText) {
    super(`Invalid response ${status}: ${statusText}`);
    this.status = status;
  }
}

export default ApiClient;
