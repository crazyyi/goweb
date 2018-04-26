import React, { Component } from 'react';
import './../css/app.css';
import DataForm from './DataForm';

class PersonItem extends Component {
  handleClick(e) {
    e.preventDefault();
    axios.delete('/delete', {params: {id: this.props.id }})
      .then(response => {
        this.props.handler();
      })
      .catch(err => {
        console.log(err);
      })
  }

  render() {
    return (
      <tr>
        <td>{this.props.id}</td>
        <td>{this.props.first}</td>
        <td>{this.props.last}</td>
        <td><input type="button" value="delete" onClick={(e) => this.handleClick(e)}/></td>
      </tr>
    )
  }
}

class PeopleList extends Component {
  constructor(props) {
    super(props)
    this.state = { people: [] };
    this.updateHandler = this.updateHandler.bind(this);
  }

  componentDidMount() {
    this.updateHandler();
  }

  updateHandler() {
    this.serverRequest = 
      axios
        .get("/people")
        .then((result) => {
          this.setState({ people: result.data })
        })
  }

  render() {
    const people = this.state.people.map((person, i) => {
      return (
        <PersonItem key={i} id={person.Id} first={person.First} last={person.Last} handler={ this.updateHandler }/>
      );
    }); 

    return (
      <div className="main">
        <span>Data:</span>
        <table>
          <tbody>
            <tr><th>Id</th><th>First</th><th>Last</th></tr>
            {people}
          </tbody>
        </table>
        <DataForm handler={ this.updateHandler }/>
      </div>
    )
  }
}

export default PeopleList;