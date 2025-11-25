document.getElementById('simForm').addEventListener('submit', async function (e) {
    e.preventDefault();

    const runBtn = document.getElementById('runBtn');
    const originalText = runBtn.innerText;
    runBtn.innerText = 'Running...';
    runBtn.disabled = true;

    const data = {
        adam: parseInt(document.getElementById('adam').value),
        eve: parseInt(document.getElementById('eve').value),
        lilith: parseInt(document.getElementById('lilith').value),
        diana: parseInt(document.getElementById('diana').value),
        eveBirth: parseFloat(document.getElementById('eveBirth').value),
        lilithBirth: parseFloat(document.getElementById('lilithBirth').value),
        dianaBirth: parseFloat(document.getElementById('dianaBirth').value),
        viableY: document.getElementById('viableY').checked,
        maxPop: parseInt(document.getElementById('maxPop').value),
        generations: parseInt(document.getElementById('generations').value),
        timelines: parseInt(document.getElementById('timelines').value)
    };

    try {
        const response = await fetch('/api/simulate', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json'
            },
            body: JSON.stringify(data)
        });

        if (!response.ok) {
            throw new Error('Simulation failed');
        }

        const result = await response.json();

        const resultsContainer = document.getElementById('results');
        const resultText = document.getElementById('resultText');

        resultText.innerText = result.summary;
        resultsContainer.classList.remove('hidden');

        // Scroll to results
        resultsContainer.scrollIntoView({ behavior: 'smooth' });

    } catch (error) {
        alert('An error occurred: ' + error.message);
    } finally {
        runBtn.innerText = originalText;
        runBtn.disabled = false;
    }
});
