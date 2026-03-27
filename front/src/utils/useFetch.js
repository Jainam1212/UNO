export const customFetch = async (url, method, body) => {
  const options = {
    method: method.toUpperCase(),
    headers: {
      "Content-Type": "application/json",
    },
    credentials: "include",
  };
  if (body) {
    options.body = JSON.stringify(body);
  }
  if (options.method === "GET") {
    delete options.body;
  }

  try {
    const response = await fetch(url, options);
    if (!response.ok) {
      const errorData = await response
        .json()
        .catch(() => ({ message: "Unknown error" }));
      throw new Error(
        `HTTP error! Status: ${response.status}. Message: ${errorData.message || response.statusText}`,
      );
    }
    return await response.json();
  } catch (error) {
    console.error("Fetch error:", error.message);
    throw error;
  }
};
