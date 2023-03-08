import { useEffect, useState } from "react";
import { Navigate } from "react-router-dom";

function Auth(props) {
    const url = process.env.REACT_APP_API_URL + "users/authtest"
    const [cred, setCred] = useState(null)
    const auth = localStorage.getItem("Authorization")
    useEffect(()=>{
        fetch(url, {method: "GET", headers:{"Authorization": auth}}).then(res=>{
            if(res.ok){
                return setCred(true);
            }
            return setCred(false);
            
        }).catch(e => {
            setCred(false)
            console.log("Auth: " + localStorage.getItem("Authorization"))
        })
    }, [cred])

    if(cred == null){
        return (
            <h1>Checking Credentials</h1>
        );
    } else if (cred == true) {
        return props.children;
    } else return <Navigate to="/signin" />
}

export default Auth