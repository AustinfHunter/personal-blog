import { useState } from "react"

function SignInForm(props) {
    return(
    <div className="signin-form">
        <h2>Sign In</h2>
        <form className="form" method="post" onSubmit={props.handleSubmit}>
            <input type="email" name="email" placeholder="Email"/>
            <input type="password" name="password" placeholder="Password"/>
            <button className="submit" type="submit">Submit</button>
        </form>
    </div>
    )
}

export default SignInForm