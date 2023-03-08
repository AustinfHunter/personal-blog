import { useEffect, useState } from "react";
import { Navigate } from "react-router-dom";
import SignInForm from "../components/SignInForm";

function SignIn () {
    const [isSignedIn, setIsSignedIn] = useState(false)


    const handleSubmit = (e) => {
        e.preventDefault();
        const form = e.target;
        const formData = new FormData(form)
        const formJSON = JSON.stringify(Object.fromEntries(formData.entries()))
        const signinURL = process.env.REACT_APP_API_URL + "users/signin"
        fetch(signinURL, {method:"POST", body: formJSON}).then(res => {
            if(res.ok){
            localStorage.setItem("Authorization", res.headers.get("Authorization"));
            localStorage.setItem("UID", res.headers.get("Uid"));
            setIsSignedIn(true)
            }      
        })
    }
    

    if(!isSignedIn){
        return(
    <div className="signin">
        <SignInForm handleSubmit={handleSubmit}/>
    </div>
    )} else return <Navigate to="/"></Navigate>
}

export default SignIn