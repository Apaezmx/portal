document.addEventListener('DOMContentLoaded', () => {
    const searchForm = document.getElementById('search-form');
    const searchInput = document.getElementById('search-input');
    const resultsContainer = document.getElementById('results-container');

    searchForm.addEventListener('submit', async (e) => {
        e.preventDefault();
        const query = searchInput.value.trim();
        if (!query) return;

        resultsContainer.innerHTML = '<p>Searching...</p>';

        try {
            const response = await fetch('/search', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                },
                body: JSON.stringify({ query: query }),
            });

            if (!response.ok) {
                throw new Error(`HTTP error! status: ${response.status}`);
            }

            const data = await response.json();
            displayResults(data);
        } catch (error) {
            console.error('Error fetching search results:', error);
            resultsContainer.innerHTML = '<p>Error fetching results. Please try again.</p>';
        }
    });

    function displayResults(data) {
        resultsContainer.innerHTML = '';

        if (!data || (!data.summary && (!data.sources || data.sources.length === 0))) {
            resultsContainer.innerHTML = '<p>No results found.</p>';
            return;
        }

        if (data.summary) {
            const summaryEl = document.createElement('div');
            summaryEl.className = 'result-item';
            summaryEl.innerHTML = `<h3>Summary</h3><p>${escapeHTML(data.summary)}</p>`;
            resultsContainer.appendChild(summaryEl);
        }

        if (data.sources && data.sources.length > 0) {
            data.sources.forEach(source => {
                const sourceEl = document.createElement('div');
                sourceEl.className = 'result-item';
                sourceEl.innerHTML = `
                    <h3><a href="${escapeHTML(source.url)}" target="_blank">${escapeHTML(source.title || source.url)}</a></h3>
                    <p>${escapeHTML(source.snippet)}</p>
                `;
                resultsContainer.appendChild(sourceEl);
            });
        }
    }

    function escapeHTML(str) {
        const p = document.createElement('p');
        p.appendChild(document.createTextNode(str));
        return p.innerHTML;
    }
});
