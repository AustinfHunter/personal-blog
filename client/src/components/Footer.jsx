import github from '../img/social/github-64.png'
import linkedin from '../img/social/linkedin-64.png'
import youtube from '../img/social/youtube-64.png'

function Footer(){
    return (
        <div className="footer">
            <a href='https://www.github.com/AustinfHunter' target="_blank">
                <img src={github} />
            </a>

            <a href='https://www.linkedin.com/in/austinfhunter' target="_blank">
                <img src={linkedin} />
            </a>
        </div>
    )
}

export default Footer