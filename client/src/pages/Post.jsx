import '../styles/blog.scss'

import { useEffect, useState } from "react"
import { useParams } from "react-router-dom"
import ReactMarkdown from 'react-markdown'
import remarkGfm from 'remark-gfm'
import { Prism as SyntaxHighlighter } from 'react-syntax-highlighter'
import { coldarkDark } from 'react-syntax-highlighter/dist/esm/styles/prism'
import rehypeRaw from 'rehype-raw'
import Loading from '../components/Loading'

function Post() {
    const url = process.env.REACT_APP_API_URL + "posts/" + useParams().post
    console.log("URL: " + url)
    const [post, setPost] = useState(null)
    const [author, setAuthor] = useState("")
    const [date, setDate] = useState(new Date())

    useEffect(() => {
        fetch(url).then(res => res.json())
            .then(data => {
                setPost(data.post)
                setAuthor(data.author.FirstName + " " + data.author.LastName)
                document.title = data.post.title
                setDate(new Date(data.post.uploadDate.Time))
            })
    }, [])

    if (!post) return <Loading/>
    return (
        <div className="post">
            <div className="content">
                <h1 className='title'>{post.title}</h1>
                <div className="info">
                    <span>{author}</span>
                    <span>{date.getUTCMonth()+1}/{date.getUTCDate()}/{date.getUTCFullYear()}</span>
                </div>
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