document.addEventListener('DOMContentLoaded', function() {
    document.getElementById('callApiButton').addEventListener('click', function() {
        // The variable to pass to the BFF
        const endpointVar = 'helloworld';

        // Making a request to the BFF
        fetch(`https://admin.ryanschnabel.com/bff/${endpointVar}`)
            .then(response => response.text())
            .then(data => {
                // Handle the response data
                console.log(data);
            })
            .catch(error => {
                // Handle any errors
                console.error('Error:', error);
            });
    });
});
