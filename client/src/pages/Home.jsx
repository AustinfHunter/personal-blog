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
                    My name is Austin Hunter and this is my developer blog. This is where I'll be writing about technical topics that I find interesting, including software
                    engineering, and web development. I'll also be documenting personal projects with the goal of giving others the opportunity to learn or follow
                    along with me as I explore new programming languages, frameworks, and technical problems. 
                </p>

                <p>
                    You can find me on my <a href='https://www.linkedin.com/in/austinfhunter' target="_blank">LinkedIn</a> and <a href='https://www.github.com/AustinfHunter' target="_blank">Github</a>.
                </p>

                <p>
                    This blog is still a work in progress, so things will likely change a lot in the near future. 
                </p>
            </div>
        </div>
    )
}

export default Home