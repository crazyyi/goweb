import React, { Component } from 'react';
import layout from './../css/app.css';
import DataForm from './DataForm';
import FontAwesomeIcon from '@fortawesome/react-fontawesome'
import * as Icons from '@fortawesome/fontawesome-free-solid';

class PersonItem extends Component {
  constructor(props) {
    super(props);
    this.state = {update: false, editId: 0, firstEdit: this.props.first, lastEdit: this.props.last};
  }

  handleChange(e) {
    this.setState({
      [e.target.name]:e.target.value,
    })
  }
  
  handleClick(e) {
    e.preventDefault();
    axios.delete('/delete', {params: {id: this.props.id }})
        .then(response => {
          this.props.handler(this.props.id);
        })
        .catch(err => {
          console.log(err);
        });
  }

  handleIconClick(e) {
    switch(e.currentTarget.id) {
      case "edit": {
        if (!this.props.updating) {
          this.props.setUpdateStatus(true);
          if (!this.state.update) {
            this.setState({update: true, editId: this.props.id});
          }
        }
        break;
      }
      case "cancel": {
        if (this.props.updating) {
          this.props.setUpdateStatus(false);
          if (this.state.update) {
            this.setState({update: false});
          }
        }

        break;
      }
      case "check": {
        const { id, first, last } = this.props;
        const { firstEdit, lastEdit } = this.state;
        // Update existing record
        if (firstEdit !== first || lastEdit !== last) {
            console.log('post');
            axios.post('/update', {
              id: id,
              firstname: this.state.firstEdit,
              lastname: this.state.lastEdit
            }, {
              headers: {
                'Content-Type': 'application/json',
              }
            })
            .then(response => {
              this.props.setUpdateStatus(false);
              this.setState({update: false});
            })
            .catch(err => {
              console.log(err);
            });
        } else {
          this.props.setUpdateStatus(false);
          this.setState({update: false});
        }
        
        break;
      }
      default: break;
    }
  }

  render() {
    const { update, firstEdit, lastEdit } = this.state;

    const editStatus = update? Icons.faTimes : Icons.faEdit;

    const editId = update? "cancel" : "edit";

    return (
      <tr>
        <td>{this.props.id}</td>
        <td>{update?<input type="text" className={`${layout.editInput}`} name="firstEdit" value={firstEdit} onChange={e => this.handleChange(e)} /> : this.state.firstEdit}</td>
        <td>{update?<input type="text" className={`${layout.editInput}`} name="lastEdit" value={lastEdit} onChange={e => this.handleChange(e)}/> : this.state.lastEdit}</td>
        <td>
          <div className={`${layout.editDiv}`} name="edit" >
            <FontAwesomeIcon id={ editId } icon={editStatus} onClick={e => this.handleIconClick(e)}/>
            { update? <FontAwesomeIcon id="check" className={`${layout.checkIcon}`} icon={Icons.faCheck} onClick={e => this.handleIconClick(e)}/> : null }
          </div>
        </td>
        <td>
          <div className={`${layout.deleteDiv}`} name="delete" onClick={e => this.handleClick(e)}> 
            <FontAwesomeIcon icon={Icons.faTrashAlt} />
          </div>
        </td>
      </tr>
    )
  }
}

class PeopleList extends Component {
  constructor(props) {
    super(props)
    this.state = { people: [], updating: false };
    this.updateHandler = this.updateHandler.bind(this);
    this.updateStatus = this.updateStatus.bind(this);
  }

  componentDidMount() {
    this.updateHandler();
  }

  updateStatus(v) {
    this.setState({
      updating: v
    });
  }

  updateHandler() {
    this.serverRequest = 
      axios
        .get("/people")
        .then((result) => {
          this.setState({ people: result.data })
        })
  }

  refreshData(id) {
    this.setState((prevState, props) => {
      return {
        people: prevState.people.filter(person => person.Id !== id)
      }
    })
  }

  render() {
    const people = this.state.people.map((person, i) => {
      return (
        <PersonItem key={person.Id} id={person.Id} first={person.First} last={person.Last} 
          updating={this.state.updating} setUpdateStatus={ this.updateStatus } handler={ (id) => this.refreshData(id) }/>
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