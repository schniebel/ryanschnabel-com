document.addEventListener('DOMContentLoaded', function() {

    // Function to add a user
    function addUser() {
        const inputText = document.getElementById('inputText');
        const requestData = {
            endpointVar: 'addUser',
            inputText: inputText.value
        };

        fetch('https://admin.ryanschnabel.com/bff/', {
            method: 'PUT',
            headers: {
                'Content-Type': 'application/json'
            },
            body: JSON.stringify(requestData)
        })
        .then(response => response.text())
        .then(data => {
            console.log(data);
            fetchAndDisplayUsers(); 
            inputText.value = '';
        })
        .catch(error => {
            console.error('Error:', error);
        });
    }

    document.getElementById('addUserButton').addEventListener('click', addUser);

    document.getElementById('inputText').addEventListener('keypress', function(event) {
        if (event.key === 'Enter') {
            event.preventDefault();
            addUser();
        }
    });

    fetchAndDisplayUsers();
});

function fetchAndDisplayUsers() {
    fetch('https://admin.ryanschnabel.com/bff/getUsers', {
        method: 'GET',
        headers: {
            'Content-Type': 'application/json'
        }
    })
    .then(response => response.text())
    .then(data => {
        const usersArray = data.split(',');

        const usersDiv = document.getElementById('whitelistedUsers');
        usersDiv.innerHTML = '';

        usersArray.forEach(user => {
            user = user.trim();
            if (user) {
                const userSpan = document.createElement('span');
                userSpan.textContent = user;
                usersDiv.appendChild(userSpan);

                if (user !== 'ryan.d.schnabel@gmail.com' && user !== 'ryan.schnabel@gmail.com') {
                    const removeButton = document.createElement('button');
                    removeButton.textContent = 'Remove User';
                    removeButton.onclick = function() { removeUser(user); };
                    usersDiv.appendChild(removeButton);
                }

                usersDiv.appendChild(document.createElement('br'));
            }
        });
    })
    .catch(error => {
        console.error('Error:', error);
    });
}

function removeUser(userEmail) {
    const requestData = {
        endpointVar: 'removeUser',
        inputText: userEmail
    };

    fetch('https://admin.ryanschnabel.com/bff/', {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json'
        },
        body: JSON.stringify(requestData)
    })
    .then(response => response.text())
    .then(data => {
        console.log(data);
        fetchAndDisplayUsers(); // Refresh the user list after removing a user
    })
    .catch(error => {
        console.error('Error:', error);
    });
}
