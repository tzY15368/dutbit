import React,{Component} from 'react'
import { Modal,Button,Descriptions, Badge, Input, Space} from 'antd';
import { EyeInvisibleOutlined,EyeTwoTone,CloseOutlined } from '@ant-design/icons';
import './App.css'
import {EditOutlined} from "@ant-design/icons";
export default class App extends Component {
    constructor(props) {
        super(props);
        this.state = {
            userCanAmend:["username","password","email"],
            userInfo:{},
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
        /*fetch("https://www.dutbit.com/userservice/v1/userinfo",{
            method:"GET",
            headers:{
                'Content-Type': 'application/json',
                'Cookies':'SESSIONID=e0b01ee841efa57b3c8e316ca27139f5'
            }
        }).then(res=>{
            if(res.status!==200){
                return Promise.reject(res)
            }
            return res.json()
        }).then(res=>{
            console.log(res)
            this.setState({userInfo:res})
        }).catch(err=>{
            console.log(`error:${err}`)
        })*/
        let jinfo = {"_id":"5f6ea1df06f16c9c47adbaa7","confirmation":"map[]","created_at":"1601085919856","email":"tzy15368@outlook.com","ip":"111.117.123.72","last_login_ip":"111.117.123.72","last_login_time":"1602759214294","password":"bf278df12620a00e3e76a8a9cce6f705","role":"[]","site":"{\"super_admin\":true}","username":"LN"}
        let roles = {"vol_time_admin":"志愿时长查询管理员","super_admin":"Super Admin"}
        let userInfoAmend = {}
        for(let i=0;i<this.state.userCanAmend.length;i++){
            userInfoAmend[this.state.userCanAmend[i]] = jinfo[this.state.userCanAmend[i]]
        }
        this.setState({
            userInfo:jinfo,
            roles:roles,
            userInfoAmend:userInfoAmend
        })
    };
    showModal = ()=>{
        this.setState({
            modalVisible: true,
        });
    };
    handleModalCancel = e => {
        console.log(e);
        this.setState({
            modalVisible: false,
        });
    };
    handleModalOk = e => {
        console.log(e);
        this.setState({
            modalLoading:true
        });
        setTimeout(()=>{this.setState({modalLoading:false,modalVisible:false})},2000)
    };
    handleAmend = (e)=>{
        let old_val = this.state.userInfoAmend
        old_val[e.target.placeholder] = e.target.value
        console.log(old_val)
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
    }
    render(){
        const {userInfo,roles,modalVisible,userInfoAmend,modalLoading,userPasswordAmend,oldPasswordInput,newPasswordInput} = this.state;
        const siteInfo = JSON.parse(userInfo.site);
        console.log(userInfoAmend);
        console.log('old psw val:',oldPasswordInput,'new psw val',newPasswordInput)
        return (
            <div>
                <Modal
                    title="Basic Modal"
                    visible={modalVisible}
                    onOk={this.handleModalOk}
                    onCancel={this.handleModalCancel}
                    footer={[
                        <Button key="back" onClick={this.handleModalCancel}>
                            Return
                        </Button>,
                        <Button key="submit" type="primary" loading={modalLoading} onClick={this.handleModalOk}>
                            Submit
                        </Button>,
                    ]}
                >

                    <Space direction="vertical">
                    {
                    Object.keys(userInfoAmend).map((key,value)=>{
                        if(key!=='password'){
                            return (
                                <div>
                                    <Input
                                        placeholder={key}
                                        addonAfter={<EditOutlined/>}
                                        value={Object.values(userInfoAmend)[value]}
                                        onChange={this.handleAmend}/>
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
                    <Descriptions.Item label="Previous Login Time">{userInfo.last_login_time}</Descriptions.Item>
                    <Descriptions.Item label="Password">*</Descriptions.Item>
                    <Descriptions.Item label="Validation">{userInfo.confirmation}</Descriptions.Item>
                    <Descriptions.Item label="Previous Login Ip">{userInfo.ip}</Descriptions.Item>
                    <Descriptions.Item label="Role">
                        {
                            Object.keys(siteInfo).map((key,value)=>{
                                return (
                                    <div>
                                        <Badge status="processing" text="" />{roles[key]}
                                    </div>
                                )
                            })
                        }

                    </Descriptions.Item>
                </Descriptions>
            </div>
        )
    }
}
