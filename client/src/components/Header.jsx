import "../App.scss";
import Hamburger from '../img/hamburger-white.png';
import {NavLink} from 'react-router-dom';
import { useState } from "react";

function Header(props) {
    let [isExpanded, setIsExpanded] = useState(false);

    function toggle() {
        setIsExpanded(!isExpanded)
    }

    function closeMenu() {
        setIsExpanded(false)
    }


    let links = props.links.map((l)=>{
        return <NavLink onClick={closeMenu} className={({isActive})=> isActive ? "nav-link-active" : "nav-link"} key={l.name} to={l.url}>{l.name}</NavLink>
    })

    return(
        <>
        <div className="header">
            <NavLink className="navLogo" to="/">{props.logoContent}</NavLink>
            <nav className="nav-links">
                {links}
            </nav>
            <img onClick={toggle} src={Hamburger} className="hamburger"/>
        </div>
        <div className={isExpanded ? "nav-menu-expanded" : "nav-menu"}>
                {links}
        </div>
        </>
    )
}

export default Header