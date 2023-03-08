import { useEffect } from 'react'
import '../styles/home.scss'

function Home(props) {
    useEffect(()=>{
        document.title = "Austin Hunter - Home"
    }, [])

    return (
        <div className="home">
            <div className="bio">
                <h1>Hi there!</h1>
                <p>
                    My name is Austin Hunter and this is my developer blog. This is were I write about technical topics that I find interesting, including software
                    engineering, and web development. I'll also be documenting personal projects here as I work on them with the goal of allowing others to learn or follow
                    along with me as I explore new projects, frameworks, programming languages, and technical problems. 
                </p>

                <p>
                    You can find me on my <a href='https://www.linkedin.com/in/austinfhunter' target="_blank">LinkedIn</a> and <a href='https://www.github.com/AustinfHunter' target="_blank">Github</a>.
                </p>

                <p>
                    This is version 1.0 of this blog, so it is likely that things will change a lot in the near future. Some planned changes include support for comments on posts,
                    support for a wider range of media in posts, and an RSS Feed for the blog. 
                </p>
            </div>
        </div>
    )
}

export default Home