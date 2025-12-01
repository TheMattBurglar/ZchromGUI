document.getElementById('simForm').addEventListener('submit', async function (e) {
    e.preventDefault();

    const runBtn = document.getElementById('runBtn');
    const originalText = runBtn.innerText;
    runBtn.innerText = 'Initializing...';
    runBtn.disabled = true;

    const resultsContainer = document.getElementById('results');
    const resultText = document.getElementById('resultText');

    // Hide previous results
    resultsContainer.classList.add('hidden');
    resultText.innerText = "";

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
        if (typeof initSimulation !== 'function') {
            throw new Error('WASM module not loaded yet. Please wait a moment.');
        }

        // Initialize simulation
        const initErr = initSimulation(JSON.stringify(data));
        if (initErr) {
            throw new Error(initErr);
        }

        // Run in batches
        const batchSize = 10; // Adjust based on performance needs

        function runLoop() {
            try {
                const status = runBatch(batchSize);

                if (status.finished) {
                    runBtn.innerText = 'Finalizing...';
                    const resultJson = getResults();
                    if (resultJson.startsWith("Error")) {
                        throw new Error(resultJson);
                    }
                    const result = JSON.parse(resultJson);

                    resultText.innerText = result.summary;
                    resultsContainer.classList.remove('hidden');
                    resultsContainer.scrollIntoView({ behavior: 'smooth' });

                    runBtn.innerText = originalText;
                    runBtn.disabled = false;
                } else {
                    runBtn.innerText = `Running Timeline ${status.completed + 1} / ${status.total}...`;
                    // Yield to UI thread
                    setTimeout(runLoop, 0);
                }
            } catch (err) {
                alert('An error occurred during simulation: ' + err.message);
                runBtn.innerText = originalText;
                runBtn.disabled = false;
            }
        }

        // Start the loop
        runLoop();

    } catch (error) {
        alert('An error occurred: ' + error.message);
        runBtn.innerText = originalText;
        runBtn.disabled = false;
    }
});
