async function sendQuery() {
    const query = document.getElementById("query").value;
    const responseElement = document.getElementById("response");
  
    if (!query.trim()) {
      responseElement.textContent = "Please enter a query.";
      return;
    }
  
    try {
      const res = await fetch("http://localhost:8081/api/query", { // لاحظ التغيير هنا
        method: "POST",
        headers: {
          "Content-Type": "application/json"
        },
        body: JSON.stringify({ query: query })
      });
  
      const data = await res.text();
      responseElement.textContent = `Master Response:\n${data}`;
    } catch (err) {
      responseElement.textContent = `Error:\n${err}`;
    }
  }
  