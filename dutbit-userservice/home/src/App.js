import React,{Component} from 'react'
import { Modal,Button,Descriptions, Badge, Input, Space,message} from 'antd';
import { EyeInvisibleOutlined,EyeTwoTone,CloseOutlined } from '@ant-design/icons';
import './App.css'
import {EditOutlined} from "@ant-design/icons";
import CONFIG from "./config";
import js_date_time from "./utils/datetime";
import ValidateEmail from "./utils/ValidateEmail";
export default class App extends Component {
    constructor(props) {
        super(props);
        this.state = {
            userCanAmend:["username","password","email"],
            userInfo:{},
            siteInfo:"{}",
            userInfoAmend:{},
            userPasswordAmend:false,
            roles:{},
            modalVisible:false,
            modalLoading:false,
            oldPasswordInput:'',
            newPasswordInput:'',
        }
    }
    componentWillMount = () => {
        console.log("mounting");
        fetch(CONFIG['USERINFO_API'],{
            method:"GET",
            headers:{
                'Content-Type': 'application/json',
                redirect: 'follow'
            }
        }).then(res=>{
            console.log(res)
            if(!res.ok){
                return Promise.reject(res.status)
            }
            if (res.redirected) {
                window.location.href = res.url;
                return Promise.reject("Did you log in?")
            }
            return res.json()
        }).then(res=>{
            console.log(res);
            let roles = {"vol_time_admin":"志愿时长查询管理员","super_admin":"Super Admin"};
            let userInfoAmend = {};
            for(let i=0;i<this.state.userCanAmend.length;i++){
                userInfoAmend[this.state.userCanAmend[i]] = {
                    value: res[this.state.userCanAmend[i]],
                    disabled:true
                }
            }
            this.setState({
                userInfo:res,
                roles:roles,
                userInfoAmend:userInfoAmend,
                siteInfo:res.site,
            })
        }).catch(err=>{
            this.ErrorMsg(`Error encountered: ${err}`);
            console.log(`Error encountered: ${err}`)
        });
        //let jinfo = {"_id":"5f6ea1df06f16c9c47adbaa7","confirmation":"map[]","created_at":"1601085919856","email":"tzy15368@outlook.com","ip":"111.117.123.72","last_login_ip":"111.117.123.72","last_login_time":"1602759214294","password":"bf278df12620a00e3e76a8a9cce6f705","role":"[]","site":"{\"super_admin\":true}","username":"LN"}
    };
    showModal = ()=>{
        this.setState({
            modalVisible: true,
        });
    };
    handleModalCancel = () => {
        this.setState({
            modalVisible: false,
        });
    };
    SuccessMsg = (msg) => {
        message.success(msg);
    };

    ErrorMsg = (err) => {
        message.error(err);
    };
    handleModalOk = e => {
        let result = {};
        for(let i=0;i<this.state.userCanAmend.length;i++){
            if(this.state.userCanAmend[i]==="password"){
                continue
            }
            result[this.state.userCanAmend[i]] = this.state.userInfoAmend[this.state.userCanAmend[i]]["value"]
        }
        result['old_password'] = this.state.oldPasswordInput;
        result['new_password'] = this.state.newPasswordInput;
        console.log("update result: ",result);
        const min_length = CONFIG['MINIMAL_PASSWORD_LENGTH'];
        if(result['username']===''){
            this.ErrorMsg(`Invalid username`);
            return;
        }
        if(!ValidateEmail(result['email'])){
            this.ErrorMsg(`Invalid email address`);
            return;
        }
        if(
            result['old_password']!==''&&result['new_password']!==''&&
            (result['old_password'].length<min_length||result['new_password'].length<min_length)
        ){
            this.ErrorMsg(`Password is shorter than ${min_length} characters`);
            return;
        }
        this.setState({
            modalLoading:true
        });
        fetch(CONFIG["USERINFO_API"],{
            method:"PUT",
            headers:{
                'Content-Type': 'application/json',
                redirect: 'follow'
            },
            body: JSON.stringify(result)
        }).then(res=>{
            if (res.redirected) {
                window.location.href = res.url;
                return Promise.reject("Did you log in?")
            }
            this.setState({modalLoading:false});
            if(!res.ok){
                return Promise.reject(res.status)
            }
            return res.json()
        }).then((res)=>{
            console.log(res);
            if(res.success){
                this.SuccessMsg(res.details);
                this.setState({modalVisible:false});
            } else {
                return Promise.reject(res.details)
            }
        }).then(this.componentWillMount).catch(err=>{
            console.log(err);
            this.ErrorMsg(`Error encountered: ${err}`);
        });
        setTimeout(()=>{this.setState({modalLoading:false,modalVisible:false})},2000)
    };
    handleAmend = (e)=>{
        let old_val = this.state.userInfoAmend;
        old_val[e.target.placeholder]["value"] = e.target.value;
        //console.log(old_val)
        this.setState({
            userInfoAmend:old_val
        })
    };
    handlePswChange = (e)=>{
        if(e.target.placeholder==="new password"){
            this.setState({newPasswordInput:e.target.value})
        } else {
            this.setState({oldPasswordInput:e.target.value})
        }
    };
    togglePasswordAmend = (e)=>{
        this.setState({
            userPasswordAmend:!this.state.userPasswordAmend
        })
    };
    toggleDisabled = (e,r)=>{
        //console.log(e,r)
        //e.target.enable()
        let old_state = this.state.userInfoAmend;
        old_state[e]['disabled'] = !old_state[e]['disabled'];
        this.setState({
            userInfoAmend:old_state
        })
    };
    render(){
        const {userInfo,roles,modalVisible,userInfoAmend,modalLoading,userPasswordAmend,oldPasswordInput,newPasswordInput} = this.state;
        const siteInfo = JSON.parse(this.state.siteInfo);
        return (
            <div>
                <Modal
                    title="Edit Personal Info"
                    visible={modalVisible}
                    onOk={this.handleModalOk}
                    onCancel={this.handleModalCancel}
                    footer={[
                        <Button key="back" onClick={this.handleModalCancel}>
                            Cancel
                        </Button>,
                        <Button key="submit" type="primary" loading={modalLoading} onClick={this.handleModalOk}>
                            Submit
                        </Button>,
                    ]}
                >

                    <Space direction="vertical">
                    {
                    Object.keys(userInfoAmend).map((key,index)=>{
                        if(key!=='password'){
                            return (
                                <div>
                                    <Input
                                        placeholder={key}
                                        addonAfter={
                                            <EditOutlined
                                                onClick={this.toggleDisabled.bind(this,Object.keys(userInfoAmend)[index])}
                                            />
                                        }
                                        value={Object.values(userInfoAmend)[index]["value"]}
                                        onChange={this.handleAmend}
                                        disabled={Object.values(userInfoAmend)[index]["disabled"]}
                                    />
                                </div>
                            )
                        } else {
                            if(!userPasswordAmend){
                                return (
                                        <Input.Password
                                            placeholder="input password"
                                            value={"*****************"}
                                            addonAfter={<EditOutlined onClick={this.togglePasswordAmend}/>}
                                            iconRender={visible => (visible ? <EyeTwoTone /> : <EyeInvisibleOutlined />)}
                                            disabled={true}
                                        />
                                )
                            } else {
                                return (
                                    <div>
                                        <Space direction="vertical">
                                            <Input.Password
                                                placeholder="new password"
                                                value={newPasswordInput}
                                                onChange={this.handlePswChange}
                                                addonAfter={<CloseOutlined onClick={this.togglePasswordAmend}/>}
                                            />
                                            <Input.Password
                                                value={oldPasswordInput}
                                                placeholder="old password"
                                                iconRender={visible => (visible ? <EyeTwoTone /> : <EyeInvisibleOutlined />)}
                                                onChange={this.handlePswChange}
                                                addonAfter={<CloseOutlined onClick={this.togglePasswordAmend}/>}
                                            />
                                        </Space>
                                    </div>
                                )
                            }

                        }
                    })
                }

                    </Space>
                </Modal>
                <Descriptions
                    title="Personal Info"
                    bordered
                    column={{ xxl: 4, xl: 3, lg: 3, md: 3, sm: 2, xs: 1 }}
                    extra={<EditOutlined onClick={this.showModal}/>}
                >
                    <Descriptions.Item label="Username" onMouse>{userInfo.username}</Descriptions.Item>
                    <Descriptions.Item label="Email">{userInfo.email}</Descriptions.Item>
                    <Descriptions.Item label="Previous Login Time">{js_date_time(parseInt(userInfo.last_login_time))}</Descriptions.Item>
                    <Descriptions.Item label="Previous Update Time">{js_date_time(parseInt(userInfo.updated_at))}</Descriptions.Item>
                    <Descriptions.Item label="Validation">{userInfo.confirmation}</Descriptions.Item>
                    <Descriptions.Item label="Previous Login Ip">{userInfo.ip}</Descriptions.Item>
                    <Descriptions.Item label="Role">
                        <div>
                            <Badge status="success" text="" />{'Default User'}
                            {
                                Object.keys(siteInfo).map((key,value)=>{
                                    return (
                                        <div id={value}>
                                            <Badge status="processing" text="" />{roles[key]}
                                        </div>
                                    )
                                })

                            }
                        </div>
                    </Descriptions.Item>
                </Descriptions>
            </div>
        )
    }
}
