type Props = {
    onClick: React.MouseEventHandler
}

const Arrow = ({ onClick }: Props) => {
    return (
        <svg xmlns="http://www.w3.org/2000/svg" id="Layer_1" data-name="Layer 1" viewBox="0 0 24 24" onClick={onClick}><path d="M23.73,11.356l-5.154-5.087c-.581-.574-1.575-.167-1.575,.644v3.587H1.5c-.828,0-1.5,.671-1.5,1.5s.672,1.5,1.5,1.5h15.5v3.587c0,.811,.994,1.218,1.575,.644l5.154-5.087c.36-.356,.36-.932,0-1.288Z" /></svg>

    )
}

export default Arrow