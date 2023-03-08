import '../styles/blog.scss'

import { useEffect, useState } from "react"
import { NavLink } from "react-router-dom";
import Pagination from "../components/Pagination";

function BlogCard(props) {
    const p = props.post
    const date = new Date(p.uploadDate.Time)
    return(
        <div className="blog-card">
            <NavLink  to={"/posts/" + p.slug}>
            <h2>{p.title}</h2>
            <h4>{p.author}</h4>
            </NavLink>
            <h3>{date.getUTCMonth()+1}/{date.getUTCDate()+1}/{date.getUTCFullYear()}</h3>
        </div>
    )
}

function Blog() {
    const [posts, setPosts] = useState([]);
    const [page, setPage] = useState(0)
    const [numPosts, setNumPosts] = useState(0)
    const url = process.env.REACT_APP_API_URL + "posts";

    useEffect(()=>{
        var qString = "?offset=" + page*5
        fetch(url + qString).then(res => res.json())
        .then(data => {
            setPosts(data.posts)
            setNumPosts(data.numPosts)
        })
        document.title = "Austin Hunter - Blog";
    }, [page])

    return(
        <div>
            <div className="blog">
                {posts.map(post => <BlogCard className="blog-cards" key={post.id} post={post}></BlogCard>)}
            </div>
            <Pagination parentPage={setPage} pageSize={5} numItems={numPosts}/>
        </div>
    )
}

export default Blog