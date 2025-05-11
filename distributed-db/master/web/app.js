function sendQuery() {
    const query = document.getElementById('query').value;
    fetch('/api/query', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ query: query })
    })
    .then(response => response.json())
    .then(data => {
        document.getElementById('response').textContent = JSON.stringify(data, null, 2);
    });
}
