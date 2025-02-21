package models

import (
	"context"
	"fmt"
	"working/super_task/config"

	"github.com/google/generative-ai-go/genai"
	"google.golang.org/api/option"
)

// function/method for working with sentiment analysis of comment and messages
func SentimentAnalysis(env *config.Env, message string) (map[string]interface{}, error) {
	prompt := `here you are a sentiment analyst for the comment and messages
	         :your are going to analysis the positivity and negativity of the comment
			 :your ara going to analsis the positiviy and negativity of the messages
			 :comment is positive if and only if doesnt contain abusive and insults otherwise it is negative
			 :message is positive if and only if doesnt contain abusive and insults otherwise it is negative
			 :for example take message/comment: "This is insult" as negative
			 :critics message or comment is not negative 
			 :your response has to be in the form of:
			          response:true111111 for postive  or response:false for negative
			 : no explanation,no introduction just the response
			 for example
			    message="your are working very good keep it up you need to improve" for this your response has to be
				response:true111111
			 :here is the message or comment to analysis\n
			`
	prompt += fmt.Sprintf("message or comment %s", message)
	cxt := context.Background()
	client, err := genai.NewClient(cxt, option.WithAPIKey(env.GEMINI_API))
	if err != nil {
		return map[string]interface{}{
			"response": err.Error(),
		}, err
	}
	defer client.Close()
	model := client.GenerativeModel("gemini-1.5-flash")

	response, err := model.GenerateContent(cxt, genai.Text(prompt))
	if err != nil {
		return map[string]interface{}{
			"response": err.Error(),
		}, err
	}

	result := response.Candidates[0].Content.Parts

	var processed string
	for _, part := range result {
		processed += fmt.Sprintf("%v", part)
	}
	return map[string]interface{}{
		"response": len(processed) == 21,
	}, nil

}
