document.addEventListener("DOMContentLoaded", async () => {
    try {
        const response = await fetch("/api/stocks");
        const data = await response.json();

        data.forEach(stock => {
            const row = document.querySelector(`[data-symbol="${stock.symbol}"]`);
            if (row) {
                row.querySelector(".current-price").textContent = stock.currentPrice.toFixed(2);
            }
        });
    } catch (error) {
        console.error("Error fetching stock data:", error);
    }
});
