// look I'm not a web developer okay
package main
const index = `
<html>
  <script src="https://unpkg.com/react@16/umd/react.production.min.js"></script>
  <script src="https://unpkg.com/react-dom@16/umd/react-dom.production.min.js"></script>
  <script src="https://unpkg.com/babel-standalone@6.15.0/babel.min.js"></script>
  <body>

    <div id="mydiv"></div>

    <script type="text/babel">
      class Hello extends React.Component {
        constructor() {
            super();
			this.state = {
              temp: 0
            };
        }
        componentWillMount() {
		  fetch('/temp')
            .then(resp => resp.json())
			.then(temp => this.setState({ temp: temp }));
        }
        incrTemp(e) {
          e.preventDefault();
		  console.log("clicked incr");
		  fetch('/temp', {method: "POST", body: "" + (this.state.temp + 1)})
            .then(resp => resp.json())
			.then(temp => this.setState({ temp: temp }));
        }
        decrTemp(e) {
          e.preventDefault();
		  console.log("clicked decr");
		  fetch('/temp', {method: "POST", body: "" + (this.state.temp - 1)})
            .then(resp => resp.json())
			.then(temp => this.setState({ temp: 0 + temp }));
        }
        render() {
          var color = "white"
          if (this.state.temp > 85) {
            color = "yellow"
          }
          if (this.state.temp > 90) {
            color = "red"
          }
          return <div style={{background: color}}>
            <h1>
              A Flaky Web App
            </h1>
			<h3>Adjust the temperature to cause warning or error states</h3>
			<br/>
			<br/>
            Current Temperature: {this.state.temp}<br/>
            <button onClick={(e) => this.incrTemp(e)}>+</button>
            <button onClick={(e) => this.decrTemp(e)}>-</button><br/><br/>
			Warning Temperature: 85<br/>
			Critical Temperature: 90<br/>
          </div>
        }
      }

      ReactDOM.render(<Hello />, document.getElementById('mydiv'))
    </script>
  </body>
</html>
`
