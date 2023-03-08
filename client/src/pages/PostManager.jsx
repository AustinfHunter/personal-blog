import '../styles/admin.scss'

import { useEffect, useState } from "react"
import { NavLink } from "react-router-dom"

function PostPreview(props){
    const p = props.post
    const uploadDate = new Date(p.uploadDate.Time)
    const archived = p.archived ? "Archived" : "public"
    return(
        <NavLink className="postprev" to={"/admin/editor?updating=" + p.id}>
            <div>
           <h3>{p.title}</h3>
           <span>{uploadDate.getDate()}/{uploadDate.getMonth()+1}/{uploadDate.getFullYear()}</span>
           </div>
           <span>{archived}</span>
        </NavLink>
    )
}


function PostManager(props) {
    const [posts, setPosts] = useState([])
    useEffect(()=>{
        const url = process.env.REACT_APP_API_URL + "admin/posts"
        const auth = localStorage.getItem("Authorization")
        fetch(url, {method: "GET", headers:{"Authorization": auth}}).then(res=>res.json()).then(data=>setPosts(data.posts))
    }, [])

    const previews = posts.map((p)=><PostPreview post={p}/>)
    return(
        <div className="postManager">
            {[previews]}
        </div>
    )
}

export default PostManager