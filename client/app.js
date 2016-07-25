(function (window) {
	$(document).ready(main);

	function move(o, dx, dy) {
		var x = Number(o.css("left").slice(0, -2));
		var y = Number(o.css("top").slice(0, -2));
		o.css("left", String(x + dx) + "px");
		o.css("top", String(y + dy) + "px");
	}
	function main() {
		var interval;
		var controls = {
			leftDown: false,
			rightDown: false,
			upDown: false,
			downDown: false
		};
		function doFrame() {
			var alice = $("#alice");
			var dx = 0;
			var dy = 0;
			if (controls.leftDown) dx -= 10;
			if (controls.rightDown) dx += 10;
			if (controls.upDown) dy -= 10;
			if (controls.downDown) dy += 10;
			move(alice, dx, dy);
		}
		interval = window.setInterval(doFrame, 17);

		$(document).on("keyup", function(e) {
			if (e.keyCode === 37) controls.leftDown = false;
			if (e.keyCode === 38) controls.upDown = false;
			if (e.keyCode === 39) controls.rightDown = false;
			if (e.keyCode === 40) controls.downDown = false;
		});
		$(document).on("keydown", function(e) {
			if (e.keyCode === 37) controls.leftDown = true;
			if (e.keyCode === 38) controls.upDown = true;
			if (e.keyCode === 39) controls.rightDown = true;
			if (e.keyCode === 40) controls.downDown = true;
		});
	}
})(window);
