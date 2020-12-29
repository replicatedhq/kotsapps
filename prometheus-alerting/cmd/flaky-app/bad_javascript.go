package main

// look I'm not a web developer okay
const index = `
<html>
  <head>
    <script src="https://unpkg.com/react@16/umd/react.production.min.js"></script>
    <script src="https://unpkg.com/react-dom@16/umd/react-dom.production.min.js"></script>
    <script src="https://unpkg.com/babel-standalone@6.15.0/babel.min.js"></script>
    <link rel="icon" type="image/png" href="data:image/png;base64,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", 
  </head>
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
			.then(temp => this.setState({ temp: temp }));
        }
        render() {
          var color = "white"
          if (this.state.temp > 85) {
            color = "yellow"
          }
          if (this.state.temp > 90) {
            color = "red"
          }
          return <div style={{background: color, textAlign: "center"}}>
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
