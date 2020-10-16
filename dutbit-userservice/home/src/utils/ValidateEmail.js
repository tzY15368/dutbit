export default function ValidateEmail(email) {
    let reg = /^([a-zA-Z]|[0-9])(\w|\-)+@[a-zA-Z0-9]+\.([a-zA-Z]{2,4})$/;
    //if(reg.test(email)){return true}else{return false}
    return reg.test(email)
}
