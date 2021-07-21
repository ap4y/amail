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

  threads(terms, page = 0, perPage = 50) {
    return this.request("GET", `/search/${terms}?page=${page}&per=${perPage}`);
  }

  thread(threadId) {
    return this.request("GET", `/threads/${threadId}`);
  }

  updateTags(terms, tags) {
    return this.request("PUT", "/tags", { terms, tags });
  }

  sendMessage(message) {
    return this.request("POST", "/messages", message);
  }

  replyToMessage(messageId, replyTo) {
    return this.request(
      "GET",
      `/messages/${btoa(messageId)}/reply?reply-to=${replyTo}`
    );
  }
}

export class ApiError extends Error {
  constructor(status, statusText) {
    super(`Invalid response ${status}: ${statusText}`);
    this.status = status;
  }
}

export default ApiClient;
