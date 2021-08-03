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

  async sendMessage(message) {
    const formData = new FormData();
    message.to.forEach((addr) => formData.append("to[]", addr));
    message.cc.forEach((addr) => formData.append("cc[]", addr));
    formData.append("subject", message.subject);
    formData.append("body", message.body);
    Object.keys(message.headers).forEach((key) =>
      formData.append(`headers[${key}]`, message.headers[key])
    );
    message.attachments.forEach((file) =>
      formData.append("attachments[]", file, file.name)
    );

    const res = await fetch(`${this.baseURL}/messages`, {
      method: "POST",
      body: formData,
    });

    if (res.ok) {
      return null;
    }

    throw new ApiError(res.status, res.statusText);
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
