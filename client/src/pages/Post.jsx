import '../styles/blog.scss'

import { useEffect, useState } from "react"
import { useParams } from "react-router-dom"
import ReactMarkdown from 'react-markdown'
import remarkGfm from 'remark-gfm'
import { Prism as SyntaxHighlighter } from 'react-syntax-highlighter'
import { coldarkDark } from 'react-syntax-highlighter/dist/esm/styles/prism'
import rehypeRaw from 'rehype-raw'

function Post() {
    const url = process.env.REACT_APP_API_URL + "posts/" + useParams().post
    console.log("URL: " + url)
    const [post, setPost] = useState({})

    useEffect(() => {
        fetch(url).then(res => res.json())
            .then(data => {
                setPost(data.post)
                document.title = data.post.title
            })
    }, [])

    return (
        <div className="post">
            <div className="content">
                <h1>{post.title}</h1>
                <ReactMarkdown children={post.content} remarkPlugins={[remarkGfm]} rehypePlugins={[rehypeRaw]} className="markdown" components={{
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
        </div>
    )
}

export default Post