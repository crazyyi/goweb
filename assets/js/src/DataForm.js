import React, { Component } from 'react';

class DataForm extends Component {
  constructor() {
    super();
    this.state = {firstname: '', lastname: ''};
    this.handleChange = this.handleChange.bind(this);
    this.handleSubmit = this.handleSubmit.bind(this);
  }

  handleSubmit(event) {
    event.preventDefault();
    axios.post('/create', this.state, {
        headers: {
          'Content-Type': 'application/json',
        }
      })
      .then(response => {
        console.log(response);
        this.props.handler();
      })
      .catch(err => {
        console.log(err);
      });
  }

  handleChange(e) {
    this.setState({
      [e.target.id]:e.target.value
    })
  }

  render() {
    return (
      <div className="dataform">
        <form onSubmit={this.handleSubmit}>
          <label>
            Firstname: 
          </label>
          <div>
          <input type="text" id="firstname" value={this.state.firstname} 
              onChange={this.handleChange} />
          </div>
          <br />
          <label>
            Lastname:
          </label>
          <div>
          <input type="text" id="lastname" value={this.state.lastname} 
              onChange={this.handleChange} />
          </div>
          <br />
          <input type="submit" value="Submit" />
        </form>
      </div>
    )
  }
}

export default DataForm;