window.onload = function() {
    var cards = document.getElementsByClassName("card");
    var small1 = [];
    var small2 = [];
    var small3 = [];
    for (var i = 0; i < cards.length; i++) { 
        var c = cards[i];
        var e = c.getElementsByClassName("end")[0]; 
        var cl = c.getBoundingClientRect().left; 
        var el = e.getBoundingClientRect().left;
        if (el - cl < 2) { 
            c.style.width = '63mm'; 
            small1.push(c);
        } else if (el - cl < 300) {
            c.style.width = '126mm';
            small2.push(c);
        } else if (el - cl < 500) {
            c.style.width = '189mm';
            small3.push(c);
        }
    }
    for (var i = 0; i < small3.length; i++) {
        document.body.appendChild(small3[i]);
        document.body.appendChild(small1.shift());
    }
    for (var i = 0; i < small2.length; i++) {
        document.body.appendChild(small2[i]);
    }
    for (var i = 0; i < small1.length; i++) {
        document.body.appendChild(small1[i]);
    }
};
