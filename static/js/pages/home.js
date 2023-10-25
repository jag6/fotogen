const changeHomeText = document.getElementById('home-header');
const splitScreenBackground = document.getElementById('split-screen');

const shuffleColors = (array) => {
    let currentIndex = array.length;

    while(currentIndex != 0){
        let randomIndex = Math.floor(Math.random() * array.length);
        currentIndex -=1;

        let temp = array[currentIndex];
        array[currentIndex] = array[randomIndex];
        array[randomIndex] = temp;
    }

    return array;
};
const setSplitScreenBackgroundColor = () => {
    let colors = [
        '#1A88BD', '#0A4F70', '#BDB71A', '#3D0209', '#3D151A',
        '#C84453', '#007018', '#FF243D', '#BD961A', '#573326', '#8A0615'
    ]; 
    shuffleColors(colors);
    splitScreenBackground.style.backgroundImage = `linear-gradient(to right, ${colors[0]}, ${colors[1]})`; 
    // splitScreenBackground.style.backgroundColor = colors[0]; 
}

if(changeHomeText) {
    setTimeout(() => {
        changeHomeText.querySelector('#white-text').textContent = 'YOUR';
    }, 5000);
    setTimeout(() => {
        changeHomeText.querySelector('#white-text').textContent = 'OUR';
    }, 6500);
    setTimeout(() => {
        setInterval(() => {
            setSplitScreenBackgroundColor();
        }, 1500);
    }, 6500);
}