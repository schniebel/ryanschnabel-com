document.addEventListener('DOMContentLoaded', function() {
    document.getElementById('callApiButton').addEventListener('click', function() {
        var inputText = document.getElementById('inputText').value;
        document.getElementById('apiResponse').textContent = 'Yolo ' + inputText;
    });
});