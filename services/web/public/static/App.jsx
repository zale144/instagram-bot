const {Table, Tabs, Tab, Modal, Button} = ReactBootstrap;
let token = localStorage.getItem('access_token');
let apiURL;
window.globalReactFunctions = {};

class App extends React.Component {
    constructor(props) {
        super(props);
        this.state = {
            followedData: [],
            jobs: [],
            processed: [],
            tabKey: 1,
            filterText: "",
            instaSearch: false,
            instaAccount: {},
            isSending: false,
        };
        this.processProfile = this.processProfile.bind(this);
        this.filterChanged = this.filterChanged.bind(this);
        this.handleSelect = this.handleSelect.bind(this);
        this.loadJobs = this.loadJobs.bind(this);
        window.globalReactFunctions.loadJobs = this.loadJobs;
    }

    loadFollowed() {
        fetch(apiURL + '/api/followed', {
            method: 'GET',
            headers: new Headers({
                'Authorization': 'Bearer '+ token,
            }),
        })
        .then(res => res.json())
        .then(
            (result) => {
                this.setState({
                    followedData: result,
                });
            },
            (error) => {
                // TODO handle
            }
        );
    }

    loadJobs() {
        fetch(apiURL + '/api/jobs', {
            method: 'GET',
            headers: new Headers({
                'Authorization': 'Bearer '+ token,
            }),
        })
            .then(res => res.json())
            .then(
                (result) => {
                    this.setState({
                        jobs: result,
                    });
                },
                (error) => {
                    // TODO handle
                }
            );
    }

    processProfile(user) {
        this.setState({
            followedData: this.state.followedData.map(v=>({
                ...v,
                isSending: v.Username === user
            }))
        });

        fetch(apiURL + '/api/process/' + user + '?crop-h=340&crop-w=270&height=360&width=300&title=Register%20at%20My-Site', {
            method: 'GET',
            headers: new Headers({
                'Authorization': 'Bearer '+ token,
            }),
        })
        .then(res => res.text())
        .then(
            ((result) => {
                this.loadFollowed();
            }).bind(this),
            (error) => {
                // TODO handle
            }
        )
    }

    handleSelect(tabKey) {
        let refreshProcessed = false;
        switch (tabKey) {
            case 1: this.loadFollowed();
            break;
            case 2: this.loadJobs();
            break;
            case 3: refreshProcessed = !this.state.refreshProcessed;
            break;
        }
        this.setState({
            tabKey: tabKey,
            refreshProcessed: refreshProcessed,
        });
    }

    filterChanged(e) {
        this.setState({
            filterText: e.target.value,
            instaSearch: false
        })
    }

    filterUsers() {
        return this.state.followedData.filter(v => v.Username.includes(this.state.filterText));
    }

    searchInstagram() {
        const username = document.getElementById('search').value;
        if (username === "") {
            alert('Search query cannot be empty');
            return;
        }
        fetch(apiURL + '/api/search/' + username, {
            method: 'GET',
            headers: new Headers({
                'Authorization': 'Bearer '+ token,
            }),
        })
        .then(res => res.json())
        .then(
            (result) => {
                this.setState({
                    instaSearch: true,
                    instaAccount: result
                })
            }
        );
    }

    follow() {
        this.setState({
            isSending: true
        });
        const search = document.getElementById('search');
        fetch(apiURL + '/api/follow/' + search.value, {
            method: 'GET',
            headers: new Headers({
                'Authorization': 'Bearer '+ token,
            }),
        })
        .then(res => res.json())
        .then(
            (result) => {
                search.value = '';
                this.setState({
                    instaSearch: false,
                    filterText: "",
                    isSending: false,
                })
            }
        );
    }

    componentDidMount() {
        apiURL = document.getElementById('api').value;
        this.loadFollowed();
    }

    render() {
        return <Tabs activeKey={this.state.key}
                     defaultActiveKey={1}
                     onSelect={this.handleSelect}>
                    <Tab eventKey={1} title="Followed">
                        <Followed filterChanged={this.filterChanged.bind(this)}
                              searchInstagram={this.searchInstagram.bind(this)}
                              instaSearch={this.state.instaSearch}
                              filterUsers={this.filterUsers()}
                              processProfile={this.processProfile.bind(this)}
                              instaAccount={this.state.instaAccount}
                              isSending={this.state.isSending}
                              follow={this.follow.bind(this)}/>
                    </Tab>
                    <Tab eventKey={2} title="Jobs">
                        <NewJob/>
                        <Jobs data={this.state.jobs}/>
                    </Tab>
                        <Tab eventKey={3} title="Processed users">
                            <Processed refresh={this.state.refreshProcessed}/>
                        </Tab>
                </Tabs>;

    }
}

class NewJob extends React.Component {

    constructor(props) {
        super(props);

        this.state = {
            hashtag: "",
            limit: "",
            title: ""
        };

        this.updateValue = this.updateValue.bind(this);
        this.submitForm = this.submitForm.bind(this);
    }

    updateValue(e) {
        const element = e.target;
        this.setState((prevState) => {
            let newState = {...prevState};
            newState[element.name] = element.value;
            return newState;
        });
    }

    submitForm() {
        const {hashtag, limit, title} = this.state;
        if (hashtag === "") {
            alert('hashtag cannot be empty');
            return;
        }
        if (limit === "") {
            alert('limit cannot be empty');
            return;
        }
        if (title === "") {
            alert('title cannot be empty');
            return;
        }
        fetch(apiURL + '/api/process-by-hashtag/' + hashtag + '?limit=' + limit + '&crop-h=340&crop-w=270&height=360&width=300&title=' + title, {
            method: 'GET',
            headers: new Headers({
                'Authorization': 'Bearer '+ token,
            }),
        })
            .then(res => res.text())
            .then(
                ((result) => {
                    console.log(result);
                    window.globalReactFunctions.loadJobs();
                }).bind(this),
                (error) => {
                    // TODO handle
                }
            )
    }

    render() {
        return (
            <div className="row" style={{marginTop: "1.5rem"}}>
                <div className="form-group col-md-4"><input className="form-control" placeholder="Hashtag" name="hashtag" type="text" onChange={this.updateValue} /></div>
                <div className="form-group col-md-4"><input className="form-control" placeholder="Limit" name="limit" type="number" onChange={this.updateValue} /></div>
                <div className="form-group col-md-4"><button className="btn btn-info col-md-12" onClick={this.submitForm}>Add Job</button></div>
                <div className="form-group col-md-12"><textarea className="form-control" placeholder="Message" name="title" style={{resize: "vertical", minHeight: "5em"}} onChange={this.updateValue} /></div>
            </div>
        );
    }

}

class JobModal extends React.Component {

    render() {
        const job = this.props.job;
        return (
            job && <Modal
                {...this.props}
                bsSize="large"
                aria-labelledby="contained-modal-title-lg"
            >
                <Modal.Header closeButton>
                    <Modal.Title id="contained-modal-title-lg">Current result for job #{job.ID}</Modal.Title>
                </Modal.Header>
                <Modal.Body style={{}}>
                    <Processed job={job} />
                </Modal.Body>
                <Modal.Footer>
                    <Button onClick={this.props.onHide}>Close</Button>
                </Modal.Footer>
            </Modal>
        );
    }
}

const Followed = (props) => {
    return <React.Fragment>
        <div className="input-group">
            <input id="search" onChange={props.filterChanged} placeholder="Username..." className="form-control col-md-12" style={{margin: "20px 0"}}/>
            <span className="input-group-btn">
                        <button className="btn btn-default" type="button" onClick={props.searchInstagram}>Search</button>
                    </span>
        </div>
        {!props.instaSearch && <Ul data={props.filterUsers}
                                        onClick={props.processProfile}/>}
        {
            props.instaSearch && <li className="list-group-item">
                <img src={props.instaAccount.profile_pic_url} height="45px" width="45px"/>
                <span style={{fontSize:'22px', marginLeft:'20px', marginRight:'7px'}}>
                            {props.instaAccount.username}
                        </span>
                {!props.isSending?
                    <button className="btn btn-success" onClick={props.follow}>Follow</button>
                    :
                    <img src="/static/image/ajax-loader.gif" width="25px"/>}
            </li>
        }
    </React.Fragment>
};

class Jobs extends React.Component {

    constructor(props) {
        super(props);

        this.closeModal = this.closeModal.bind(this);
        this.activateJob = this.activateJob.bind(this);

        this.state = {
            activeJob: null
        }
    }

    closeModal() {
        this.setState({
            activeJob: null
        });
    }

    activateJob(jobId) {
        const filteredJobs = this.props.data.filter(job => job.ID==jobId);
        const activeJob = filteredJobs[0] || null;
        this.setState({
            activeJob: activeJob
        });
    }

    render() {
        return (
            <React.Fragment>
                <JobModal show={!!this.state.activeJob} onHide={this.closeModal} job={this.state.activeJob}/>
                <Table responsive>
                    <thead>
                    <tr>
                        <th>ID</th>
                        <th>Hashtag</th>
                        <th>Created at</th>
                        <th>Finished at</th>
                    </tr>
                    </thead>
                    <tbody>
                    {this.props.data.map(v => {
                        return <tr>
                                <td>{v.ID}</td>
                                <td><a onClick={() => this.activateJob(v.ID)}>#{v.HashTagName}</a></td>
                                <td>{moment(v.CreatedAt).format("DD.MM.YYYY, HH:mm:ss")}</td>
                                <td>{v.FinishedAt?moment.unix(v.FinishedAt).format("DD.MM.YYYY, HH:mm:ss"):''}</td>
                            </tr>
                    })}
                    </tbody>
                </Table>
            </React.Fragment>
        );
    }
}

class Processed extends React.Component {
    constructor(props) {
        super(props);

        this.state = {
            users: [],
            page: 1,
            finished: false,
            loading: false
        };

        this.fetchUsers = this.fetchUsers.bind(this);
    }

    componentDidMount() {
        this.fetchUsers(null, false);
    }

    componentDidUpdate(prevProps, prevState, snapshot) {
        const prevJob = prevProps.job || {};
        const currentJob = this.props.job || {};

        if(JSON.stringify(prevJob) !== JSON.stringify(currentJob) && this.props.job) {
            this.setState({
                users: [],
                page: 1,
                finished: false
            }, function() {
                this.fetchUsers(null, true);
            }.bind(this));
        }
    }

    componentWillReceiveProps(props) {
        if (props.refresh) {
            this.setState({
                users: [],
                page: 1,
                finished: false
            }, function() {
                this.fetchUsers(null, true);
            }.bind(this));
        }
    }

    fetchUsers(event, rewrite = false) {
        this.setState({
            loading: true
        });
        const url = this.props.job?`/api/processed-by-job/${this.props.job.ID}/${this.state.page}`:`/api/processed/${this.state.page}`;
        fetch(apiURL + url, {
            method: 'GET',
            headers: new Headers({
                'Authorization': 'Bearer '+ token,
            }),
        })
            .then(res => res.json())
            .then(
                ((result) => {
                    this.setState({
                        users: rewrite?result:this.state.users.slice().concat(result),
                        page: this.state.page+1,
                        finished: result.length<10,
                        loading: false
                    });
                }).bind(this),
                (error) => {
                    // TODO handle
                }
            )
    }

    render() {
        return <Table responsive>
            <thead>
            <tr>
                <th>ID</th>
                <th>Username</th>
                <th>Hashtag</th>
                <th>Processed at</th>
                <th>Successful</th>
            </tr>
            </thead>
            <tbody>
            {this.state.users.map(v => {
                return <tr>
                    <td>{v.ID}</td>
                    <td>{v.Username}</td>
                    <td>#{v.Job.HashTagName}</td>
                    <td>{v.ProcessedAt?moment.unix(v.ProcessedAt).format("DD.MM.YYYY, HH:mm:ss"):''}</td>
                    <td>{v.Successful?'yes':'no'}</td>
                </tr>
            })}
            </tbody>
            {!this.state.finished && <button className="btn btn-warning" disabled={this.state.loading} onClick={this.fetchUsers}>Load more</button>}
        </Table>;
    }
}

const Ul = (props) => {
    if (!props.data) {
        return <img src="/static/image/ajax-loader.gif"/>;
    } else {
        return <ul className="list-group list-group-flush">
            {
                props.data.map((v) => {
                    const style = {fontSize:'22px', marginLeft:'20px', marginRight:'7px'};
                    let el = <a onClick={(u) => props.onClick(v.Username)} style={style}>{v.Username}</a>;
                    if (v.IsSent) {
                        el = <span><span style={style}>{v.Username}</span><img src="/static/image/sent.png" width="25px"/></span>;
                    }
                    if (v.isSending) {
                        el = <span><span style={style}>{v.Username}</span><img src="/static/image/ajax-loader.gif" width="25px"/></span>;
                    }
                    return <li className="list-group-item">
                                <img src={v.ProfilePicUrl} height="45px" width="45px"/>
                                {el}
                            </li>
                })
            }
        </ul>
    }
};

ReactDOM.render(<App />, document.getElementById('root'));
//registerServiceWorker();