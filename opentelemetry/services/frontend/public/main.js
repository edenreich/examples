(() => {
    const xhr = new XMLHttpRequest();
    xhr.setRequestHeader('Host', 'order.local');
    xhr.open('GET', 'http://localhost:8080/order');
    xhr.onload = () => {
        const data = JSON.parse(xhr.responseText);
        const el = document.querySelector("#data")
        el.innerHTML = data.message;
    }
    xhr.onerror = () => {
        console.error(xhr.statusText);
    }
    xhr.send();

})();
