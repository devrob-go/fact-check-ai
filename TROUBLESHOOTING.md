# Troubleshooting Guide

This guide helps you resolve common issues when running the Fact-Check application.

## 🐳 Docker Issues

### Frontend Build Fails
**Error**: `npm ci --only=production` fails
**Solution**: The Dockerfile has been updated to use `npm install` instead of `npm ci`. This resolves the missing `package-lock.json` issue.

### Port Already in Use
**Error**: `Bind for 0.0.0.0:8080 failed: port is already in use`
**Solution**: 
```bash
# Check what's using the port
sudo lsof -i :8080

# Kill the process or change the port in docker-compose.yml
docker-compose down
# Edit docker-compose.yml to use different ports
```

### Permission Denied
**Error**: `permission denied` when running Docker commands
**Solution**: Add your user to the docker group:
```bash
sudo usermod -aG docker $USER
# Log out and back in, or restart your system
```

## 🗄️ Database Issues

### Connection Refused
**Error**: `connection refused` to PostgreSQL
**Solution**: 
```bash
# Check if PostgreSQL container is running
docker-compose ps postgres

# View PostgreSQL logs
docker-compose logs postgres

# Restart the database
docker-compose restart postgres
```

### Database Migration Fails
**Error**: `failed to run database migrations`
**Solution**: The database will automatically create tables on first run. If issues persist:
```bash
# Reset the database
docker-compose down -v
docker-compose up -d
```

## 🔐 Authentication Issues

### Google OAuth2 Not Working
**Error**: `invalid_client` or OAuth redirect issues
**Solution**: 
1. Verify your Google OAuth2 credentials in `.env`
2. Ensure redirect URI matches exactly: `http://localhost:3000/auth/callback`
3. Check that your Google Cloud Console project has OAuth2 API enabled

### JWT Token Issues
**Error**: `invalid token` or authentication failures
**Solution**: 
1. Check JWT_SECRET in `.env` file
2. Ensure the secret is at least 32 characters long
3. Restart the backend service after changing JWT_SECRET

## 🤖 OpenAI API Issues

### API Key Invalid
**Error**: `OpenAI API returned status 401`
**Solution**: 
1. Verify your OpenAI API key in `.env`
2. Check your OpenAI account billing status
3. Ensure the API key has access to the required models

### Rate Limiting
**Error**: `rate limit exceeded`
**Solution**: 
1. Check your OpenAI API usage limits
2. Implement exponential backoff in your requests
3. Consider upgrading your OpenAI plan

## 🌐 Frontend Issues

### React App Won't Start
**Error**: `Module not found` or build failures
**Solution**: 
```bash
cd frontend
rm -rf node_modules package-lock.json
npm install
npm start
```

### API Calls Failing
**Error**: `Failed to fetch` or CORS issues
**Solution**: 
1. Ensure backend is running on port 8080
2. Check CORS configuration in backend
3. Verify API_BASE_URL in frontend config

## 🚀 Deployment Issues

### Kubernetes Deployment Fails
**Error**: `ImagePullBackOff` or pod creation issues
**Solution**: 
1. Ensure images are built and pushed to registry
2. Check Kubernetes secrets are properly configured
3. Verify resource limits and requests

### Health Checks Failing
**Error**: `liveness probe failed` or `readiness probe failed`
**Solution**: 
1. Check application logs for errors
2. Verify health check endpoints are responding
3. Adjust probe timing in Kubernetes manifests

## 🔧 General Issues

### Services Won't Start
**Error**: `exit code 1` or service failures
**Solution**: 
```bash
# View all logs
docker-compose logs

# View specific service logs
docker-compose logs backend
docker-compose logs frontend
docker-compose logs postgres

# Restart all services
docker-compose restart
```

### Memory Issues
**Error**: `out of memory` or container crashes
**Solution**: 
1. Increase Docker memory limit
2. Check resource usage: `docker stats`
3. Optimize application memory usage

### Network Issues
**Error**: `network unreachable` or connection timeouts
**Solution**: 
1. Check Docker network configuration
2. Verify firewall settings
3. Ensure ports are not blocked

## 📋 Common Commands

### Debugging
```bash
# View running containers
docker-compose ps

# View logs
docker-compose logs -f

# Execute commands in running container
docker-compose exec backend sh
docker-compose exec postgres psql -U postgres -d factcheck

# Check container resource usage
docker stats
```

### Maintenance
```bash
# Clean up unused resources
docker system prune -f

# Remove all containers and volumes
docker-compose down -v

# Rebuild images
docker-compose build --no-cache

# Update dependencies
docker-compose pull
```

## 🆘 Getting Help

If you're still experiencing issues:

1. **Check the logs**: Use `docker-compose logs` to see detailed error messages
2. **Verify configuration**: Ensure all environment variables are set correctly
3. **Check system requirements**: Ensure you have sufficient resources
4. **Search issues**: Check if your issue has been reported before
5. **Create an issue**: Provide detailed error messages and system information

## 📚 Additional Resources

- [Docker Documentation](https://docs.docker.com/)
- [Docker Compose Documentation](https://docs.docker.com/compose/)
- [Kubernetes Documentation](https://kubernetes.io/docs/)
- [Google OAuth2 Documentation](https://developers.google.com/identity/protocols/oauth2)
- [OpenAI API Documentation](https://platform.openai.com/docs)
