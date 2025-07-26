const eventSource = new EventSource('/events');

eventSource.addEventListener('message', (e) => {
    const event = JSON.parse(e.data);
    displayCommit(event);
});

function displayCommit(event) {
    const eventsDiv = document.getElementById('events');

    event.messages.forEach(message => {
        const commitDiv = document.createElement('div');
        commitDiv.className = 'commit';
        commitDiv.innerHTML = `
            <h3>${event.repository}</h3>
            <div class="commit-message">${message}</div>
            <small>Committed by ${event.committer}</small>
        `;
        eventsDiv.prepend(commitDiv);
    });
}