import '../styles/admin.scss'

import { useEffect, useState } from 'react'
import ReactMarkdown from 'react-markdown'
import remarkGfm from 'remark-gfm'
import { Prism as SyntaxHighlighter } from 'react-syntax-highlighter'
import { coldarkDark } from 'react-syntax-highlighter/dist/esm/styles/prism'
import { Navigate } from 'react-router-dom'

function PostEditor() {
    const [title, setTitle] = useState("")
    const [markdown, setMarkdown] = useState("")
    const [redirect, setRedirect] = useState(null)
    useEffect(() => {
        const sParams = new URLSearchParams(window.location.search)
        if(sParams.has("updating")){
            const url = process.env.REACT_APP_API_URL + "posts/?id=" + sParams.get("updating")
            fetch(url).then(res=>res.json()).then(data=>{
                setTitle(data.post.title);
                setMarkdown(data.post.content);
            })
        }
    }, [redirect])

    const handleTitle = (e) => {
        setTitle(e.target.value)
    }

    const handleMarkdown = (e) => {
        setMarkdown(e.target.value)
    }

    const handleSubmit = (e) => {
        e.preventDefault();
        const form = e.target;
        const formData = new FormData(form)
        const archive = formData.get("archived") === "on" ? true : false;

        const body = new Object;
        body.title = formData.get("title");
        body.authorID = Number(localStorage.getItem("UID"));
        body.content = formData.get("content");
        body.archived = archive;
        const auth = localStorage.getItem("Authorization")

        const sParams = new URLSearchParams(window.location.search)

        console.log("query" + sParams.get("updating"))

        if(sParams.has("updating")){
            body.id = Number(sParams.get("updating"))
            const url = process.env.REACT_APP_API_URL + "posts/update"
            fetch(url, { method: "POST", headers: { "Authorization": auth }, body: JSON.stringify(body) }).then(res => res.json()).then(data => {
                if(data.id){
                    setRedirect("/blog")
                }
            })
            return
        }
        const url = process.env.REACT_APP_API_URL + "posts/create"
        fetch(url, { method: "POST", headers: { "Authorization": auth }, body: JSON.stringify(body) }).then(res => res.json()).then(data => {
            if(data.id){
                setRedirect("/blog")
            }
        })
    }
    if (redirect) {
        return <Navigate to={redirect} />
    } else return (
        <div className="editor">
            <form className="editorForm" onSubmit={handleSubmit}>
                <label htmlFor="title">Title</label>
                <input type="text" id="title" name="title" className='title' value={title} onChange={handleTitle} />
                <label htmlFor="content">Content</label>
                <textarea id="content" name='content' className='content' value={markdown} onChange={handleMarkdown} />
                <button type="submit">Submit</button>
                <label htmlFor="archived">Archive?</label>
                <input type="checkbox" name="archived"></input>
            </form>
            <ReactMarkdown children={markdown} remarkPlugins={[remarkGfm]} className="markdown" components={{
                code({ node, inline, className, children, ...props }) {
                    const match = /language-(\w+)/.exec(className || '')
                    return !inline && match ? (
                        <SyntaxHighlighter
                            children={String(children).replace(/\n$/, '')}
                            style={coldarkDark}
                            language={match[1]}
                            PreTag="div"
                            {...props}
                        />
                    ) : (
                        <code className={className} {...props}>
                            {children}
                        </code>
                    )
                }
            }} />
        </div>
    )
}

export default PostEditor